package server

import (
	"errors"
	"github.com/okex/exchain/libs/cosmos-sdk/client/flags"
	"log"
	"net"
	"net/http"
	"net/rpc"
	"os"
	"reflect"
	"strings"

	"github.com/okex/exchain/libs/cosmos-sdk/server/config"
	"github.com/spf13/cobra"
)

// exchain full-node start flags
const (
	FlagListenAddr         = "rest.laddr"
	FlagExternalListenAddr = "rest.external_laddr"
	FlagUlockKey           = "rest.unlock_key"
	FlagUlockKeyHome       = "rest.unlock_key_home"
	FlagRestPathPrefix     = "rest.path_prefix"
	FlagCORS               = "cors"
	FlagMaxOpenConnections = "max-open"
	FlagHookstartInProcess = "startInProcess"
	FlagWebsocket          = "wsport"
	FlagWsMaxConnections   = "ws.max_connections"
	FlagWsSubChannelLength = "ws.sub_channel_length"

	// plugin flags
	FlagBackendEnableBackend       = "backend.enable_backend"
	FlagBackendEnableMktCompute    = "backend.enable_mkt_compute"
	FlagBackendLogSQL              = "backend.log_sql"
	FlagBackendCleanUpsTime        = "backend.clean_ups_time"
	FlagBacekendOrmEngineType      = "backend.orm_engine.engine_type"
	FlagBackendOrmEngineConnectStr = "backend.orm_engine.connect_str"

	FlagStreamEngine                        = "stream.engine"
	FlagStreamKlineQueryConnect             = "stream.klines_query_connect"
	FlagStreamWorkerId                      = "stream.worker_id"
	FlagStreamRedisScheduler                = "stream.redis_scheduler"
	FlagStreamRedisLock                     = "stream.redis_lock"
	FlagStreamLocalLockDir                  = "stream.local_lock_dir"
	FlagStreamCacheQueueCapacity            = "stream.cache_queue_capacity"
	FlagStreamMarketTopic                   = "stream.market_topic"
	FlagStreamMarketPartition               = "stream.market_partition"
	FlagStreamMarketServiceEnable           = "stream.market_service_enable"
	FlagStreamMarketNacosUrls               = "stream.market_nacos_urls"
	FlagStreamMarketNacosNamespaceId        = "stream.market_nacos_namespace_id"
	FlagStreamMarketNacosClusters           = "stream.market_nacos_clusters"
	FlagStreamMarketNacosServiceName        = "stream.market_nacos_service_name"
	FlagStreamMarketNacosGroupName          = "stream.market_nacos_group_name"
	FlagStreamMarketEurekaName              = "stream.market_eureka_name"
	FlagStreamEurekaServerUrl               = "stream.eureka_server_url"
	FlagStreamRestApplicationName           = "stream.rest_application_name"
	FlagStreamRestNacosUrls                 = "stream.rest_nacos_urls"
	FlagStreamRestNacosNamespaceId          = "stream.rest_nacos_namespace_id"
	FlagStreamPushservicePulsarPublicTopic  = "stream.pushservice_pulsar_public_topic"
	FlagStreamPushservicePulsarPrivateTopic = "stream.pushservice_pulsar_private_topic"
	FlagStreamPushservicePulsarDepthTopic   = "stream.pushservice_pulsar_depth_topic"
	FlagStreamRedisRequirePass              = "stream.redis_require_pass"

	stopServiceName  = "stop"
	stopServiceLaddr = "localhost:9000"
	stopServiceProto = "tcp"
)

const (
	// 3 seconds for default timeout commit
	defaultTimeoutCommit = 3
)

var (
	backendConf = config.DefaultConfig().BackendConfig
	streamConf  = config.DefaultConfig().StreamConfig
)

//module hook

type fnHookstartInProcess func(ctx *Context) error

type serverHookTable struct {
	hookTable map[string]interface{}
}

var gSrvHookTable = serverHookTable{make(map[string]interface{})}

func InstallHookEx(flag string, hooker fnHookstartInProcess) {
	gSrvHookTable.hookTable[flag] = hooker
}

//call hooker function
func callHooker(flag string, args ...interface{}) error {
	params := make([]interface{}, 0)
	switch flag {
	case FlagHookstartInProcess:
		{
			//none hook func, return nil
			function, ok := gSrvHookTable.hookTable[FlagHookstartInProcess]
			if !ok {
				return nil
			}
			params = append(params, args...)
			if len(params) != 1 {
				return errors.New("too many or less parameter called, want 1")
			}

			//param type check
			p1, ok := params[0].(*Context)
			if !ok {
				return errors.New("wrong param 1 type. want *Context, got" + reflect.TypeOf(params[0]).String())
			}

			//get hook function and call it
			caller := function.(fnHookstartInProcess)
			return caller(p1)
		}
	default:
		break
	}
	return nil
}

func getRealAddr(r *http.Request) net.IP {
	var	remoteIP net.IP
	// the default is the originating ip. but we try to find better options because this is almost
	// never the right IP
	if parts := strings.Split(r.RemoteAddr, ":"); len(parts) == 2 {
		remoteIP = net.ParseIP(parts[0])
	}
	// If we have a forwarded-for header, take the address from there
	if xff := strings.Trim(r.Header.Get("X-Forwarded-For"), ","); len(xff) > 0 {
		addrs := strings.Split(xff, ",")
		lastFwd := addrs[len(addrs)-1]
		if ip := net.ParseIP(lastFwd); ip != nil {
			remoteIP = ip
		}
		// parse X-Real-Ip header
	} else if xri := r.Header.Get("X-Real-Ip"); len(xri) > 0 {
		if ip := net.ParseIP(xri); ip != nil {
			remoteIP = ip
		}
	}
	return remoteIP
}

func StopServe(cleanupFunc func()) {
	srv := &http.Server{
		Addr: stopServiceLaddr,
		Handler: http.HandlerFunc( func(w http.ResponseWriter, req *http.Request) {
			// get the real IP of the user, see below
			addr := getRealAddr(req)

			// the actual vaildation - replace with whatever you want
			if addr.IsLoopback()  {
				// pass the request to the mux
				if cleanupFunc != nil {
					cleanupFunc()
					// 128 + syscall.SIGINT
					os.Exit(130)
				}
				return
			} else {
				http.Error(w, "Blocked", 401)
				return
			}
		}),
	}
	go func() {
		if err := srv.ListenAndServe(); err != http.ErrServerClosed {
			// Error starting or closing listener:
			log.Fatal("ListenAndServe error:", err)
		}
	}()
}

// StopCmd stop the node gracefully
// Tendermint.
func StopCmd(ctx *Context) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "stop",
		Short: "Stop the node gracefully",
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := rpc.DialHTTP(stopServiceProto, stopServiceLaddr)
			if err != nil {
				log.Fatal("dialing:", err)
			}

			var reply string
			err = client.Call("StopService.Stop", "stop", &reply)
			if err != nil {
				log.Fatal(err)
			}
			return nil
		},
	}
	return cmd
}

var sem *nodeSemaphore

type nodeSemaphore struct {
	done chan struct{}
}

func Stop() {
	sem.done <- struct{}{}
}

// registerRestServerFlags registers the flags required for rest server
func registerRestServerFlags(cmd *cobra.Command) *cobra.Command {
	cmd.Flags().String(FlagListenAddr, "tcp://0.0.0.0:26659", "The address for the rest-server to listen on. (0.0.0.0:0 means any interface, any port)")
	cmd.Flags().String(FlagUlockKey, "", "Select the keys to unlock on the RPC server")
	cmd.Flags().String(FlagUlockKeyHome, os.ExpandEnv("$HOME/.exchaincli"), "The keybase home path")
	cmd.Flags().String(FlagRestPathPrefix, "exchain", "Path prefix for registering rest api route.")
	cmd.Flags().String(flags.FlagKeyringBackend, flags.DefaultKeyringBackend, "Select keyring's backend (os|file|test)")
	cmd.Flags().String(FlagCORS, "", "Set the rest-server domains that can make CORS requests (* for all)")
	cmd.Flags().Int(FlagMaxOpenConnections, 1000, "The number of maximum open connections of rest-server")
	cmd.Flags().String(FlagExternalListenAddr, "127.0.0.1:26659", "Set the rest-server external ip and port, when it is launched by Docker")
	cmd.Flags().String(FlagWebsocket, "8546", "websocket port to listen to")
	cmd.Flags().Int(FlagWsMaxConnections, 20000, "the max capacity number of websocket client connections")
	cmd.Flags().Int(FlagWsSubChannelLength, 100, "the length of subscription channel")
	cmd.Flags().String(flags.FlagChainID, "", "Chain ID of tendermint node for web3")
	cmd.Flags().StringP(flags.FlagBroadcastMode, "b", flags.BroadcastSync, "Transaction broadcasting mode (sync|async|block) for web3")
	return cmd
}

// registerExChainPluginFlags registers the flags required for rest server
func registerExChainPluginFlags(cmd *cobra.Command) *cobra.Command {
	cmd.Flags().Bool(FlagBackendEnableBackend, backendConf.EnableBackend, "Enable the node's backend plugin")
	cmd.Flags().MarkHidden(FlagBackendEnableBackend)
	cmd.Flags().Bool(FlagBackendEnableMktCompute, backendConf.EnableMktCompute, "Enable kline and ticker calculating")
	cmd.Flags().MarkHidden(FlagBackendEnableMktCompute)
	cmd.Flags().Bool(FlagBackendLogSQL, backendConf.LogSQL, "Enable backend plugin logging sql feature")
	cmd.Flags().MarkHidden(FlagBackendLogSQL)
	cmd.Flags().String(FlagBackendCleanUpsTime, backendConf.CleanUpsTime, "Backend plugin`s time of cleaning up kline data")
	cmd.Flags().MarkHidden(FlagBackendCleanUpsTime)
	cmd.Flags().String(FlagBacekendOrmEngineType, backendConf.OrmEngine.EngineType, "Backend plugin`s db (mysql or sqlite3)")
	cmd.Flags().MarkHidden(FlagBacekendOrmEngineType)
	cmd.Flags().String(FlagBackendOrmEngineConnectStr, backendConf.OrmEngine.ConnectStr, "Backend plugin`s db connect address")
	cmd.Flags().MarkHidden(FlagBackendOrmEngineConnectStr)

	cmd.Flags().String(FlagStreamEngine, streamConf.Engine, "Stream plugin`s engine config")
	cmd.Flags().MarkHidden(FlagStreamEngine)
	cmd.Flags().String(FlagStreamKlineQueryConnect, streamConf.KlineQueryConnect, "Stream plugin`s kiline query connect url")
	cmd.Flags().MarkHidden(FlagStreamKlineQueryConnect)

	// distr-lock flags
	cmd.Flags().String(FlagStreamWorkerId, streamConf.WorkerId, "Stream plugin`s worker id")
	cmd.Flags().MarkHidden(FlagStreamWorkerId)
	cmd.Flags().String(FlagStreamRedisScheduler, streamConf.RedisScheduler, "Stream plugin`s redis url for scheduler job")
	cmd.Flags().MarkHidden(FlagStreamRedisScheduler)
	cmd.Flags().String(FlagStreamRedisLock, streamConf.RedisLock, "Stream plugin`s redis url for distributed lock")
	cmd.Flags().MarkHidden(FlagStreamRedisLock)
	cmd.Flags().String(FlagStreamLocalLockDir, streamConf.LocalLockDir, "Stream plugin`s local lock dir")
	cmd.Flags().MarkHidden(FlagStreamLocalLockDir)
	cmd.Flags().Int(FlagStreamCacheQueueCapacity, streamConf.CacheQueueCapacity, "Stream plugin`s cache queue capacity config")
	cmd.Flags().MarkHidden(FlagStreamCacheQueueCapacity)

	// kafka/pulsar service flags
	cmd.Flags().String(FlagStreamMarketTopic, streamConf.MarketTopic, "Stream plugin`s pulsar/kafka topic for market quotation")
	cmd.Flags().MarkHidden(FlagStreamMarketTopic)
	cmd.Flags().Int(FlagStreamMarketPartition, streamConf.MarketPartition, "Stream plugin`s pulsar/kafka partition for market quotation")
	cmd.Flags().MarkHidden(FlagStreamMarketPartition)

	// market service flags for nacos config
	cmd.Flags().Bool(FlagStreamMarketServiceEnable, streamConf.MarketServiceEnable, "Stream plugin`s market service enable config")
	cmd.Flags().MarkHidden(FlagStreamMarketServiceEnable)
	cmd.Flags().String(FlagStreamMarketNacosUrls, streamConf.MarketNacosUrls, "Stream plugin`s nacos server urls for getting market service info")
	cmd.Flags().MarkHidden(FlagStreamMarketNacosUrls)
	cmd.Flags().String(FlagStreamMarketNacosNamespaceId, streamConf.MarketNacosNamespaceId, "Stream plugin`s nacos name space id for getting market service info")
	cmd.Flags().MarkHidden(FlagStreamMarketNacosNamespaceId)
	cmd.Flags().StringArray(FlagStreamMarketNacosClusters, streamConf.MarketNacosClusters, "Stream plugin`s nacos clusters array list for getting market service info")
	cmd.Flags().MarkHidden(FlagStreamMarketNacosClusters)
	cmd.Flags().String(FlagStreamMarketNacosServiceName, streamConf.MarketNacosServiceName, "Stream plugin`s nacos service name for getting market service info")
	cmd.Flags().MarkHidden(FlagStreamMarketNacosServiceName)
	cmd.Flags().String(FlagStreamMarketNacosGroupName, streamConf.MarketNacosGroupName, "Stream plugin`s nacos group name for getting market service info")
	cmd.Flags().MarkHidden(FlagStreamMarketNacosGroupName)

	// market service flags for eureka config
	cmd.Flags().String(FlagStreamMarketEurekaName, streamConf.MarketEurekaName, "Stream plugin`s market service name in eureka")
	cmd.Flags().MarkHidden(FlagStreamMarketEurekaName)
	cmd.Flags().String(FlagStreamEurekaServerUrl, streamConf.EurekaServerUrl, "Eureka server url for discovery service of rest api")
	cmd.Flags().MarkHidden(FlagStreamEurekaServerUrl)

	// restful service flags
	cmd.Flags().String(FlagStreamRestApplicationName, streamConf.RestApplicationName, "Stream plugin`s rest application name in eureka or nacos")
	cmd.Flags().MarkHidden(FlagStreamRestApplicationName)
	cmd.Flags().String(FlagStreamRestNacosUrls, streamConf.RestNacosUrls, "Stream plugin`s nacos server urls for discovery service of rest api")
	cmd.Flags().MarkHidden(FlagStreamRestNacosUrls)
	cmd.Flags().String(FlagStreamRestNacosNamespaceId, streamConf.RestNacosNamespaceId, "Stream plugin`s nacos namepace id for discovery service of rest api")
	cmd.Flags().MarkHidden(FlagStreamRestNacosNamespaceId)

	// push service flags
	cmd.Flags().String(FlagStreamPushservicePulsarPublicTopic, streamConf.PushservicePulsarPublicTopic, "Stream plugin`s pulsar public topic of push service")
	cmd.Flags().MarkHidden(FlagStreamPushservicePulsarPublicTopic)
	cmd.Flags().String(FlagStreamPushservicePulsarPrivateTopic, streamConf.PushservicePulsarPrivateTopic, "Stream plugin`s pulsar private topic of push service")
	cmd.Flags().MarkHidden(FlagStreamPushservicePulsarPrivateTopic)
	cmd.Flags().String(FlagStreamPushservicePulsarDepthTopic, streamConf.PushservicePulsarDepthTopic, "Stream plugin`s pulsar depth topic of push service")
	cmd.Flags().MarkHidden(FlagStreamPushservicePulsarDepthTopic)
	cmd.Flags().String(FlagStreamRedisRequirePass, streamConf.RedisRequirePass, "Stream plugin`s redis require pass")
	cmd.Flags().MarkHidden(FlagStreamRedisRequirePass)
	return cmd
}
