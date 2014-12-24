package btcserver

import "github.com/hlandauf/btcnode"
import "github.com/hlandauf/btcnode/cpuminer"
import "github.com/hlandauf/btcmgmt"

type Server struct {
  Config    *Config
  node      *btcnode.Node
  rpcServer *btcmgmt.RPCServer
  cpuMiner  *cpuminer.CPUMiner
}

type Config struct {
  btcnode.NodeConfig

  RPCConfig btcmgmt.RPCServerConfig
  DisableRPC bool `long:"disablerpc" description:"Disable RPC?"`
}

func New(cfg *Config) (*Server, error) {
  s := &Server{}
  s.Config = cfg

  node, err := btcnode.NewNode(cfg.Listeners, &cfg.NodeConfig)
  if err != nil {
    return nil, err
  }

  s.node = node
  s.cpuMiner = cpuminer.New(s.node)

  if !cfg.DisableRPC {
    s.rpcServer, err = btcmgmt.NewRPCServer(cfg.RPCConfig, s)
    if err != nil {
      return nil, err
    }
  }

  return s, nil
}

func (s *Server) CPUMiner() *cpuminer.CPUMiner {
  return s.cpuMiner
}

func (s *Server) Node() *btcnode.Node {
  return s.node
}

func (s *Server) Start() {
  s.node.Start()

  if !s.Config.DisableRPC {
    s.rpcServer.Start()
  }

  if s.Config.Generate {
    s.cpuMiner.Start()
  }
}

func (s *Server) Stop() {
  s.node.Stop()

  if !s.Config.DisableRPC {
    s.rpcServer.Stop()
  }

  if s.Config.Generate {
    s.cpuMiner.Stop()
  }
}

func (s *Server) WaitForShutdown() {
  s.node.WaitForShutdown()
}
