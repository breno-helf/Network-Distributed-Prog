package commands

import (
	"errors"
	"fmt"
	"net"
	"sort"
	"time"
	"reflect"

	"../eventlog"
	"../utils"
)

//HEARTBTIME defines HeartBeat() repeat time 
const HEARTBTIME = 5 * time.Second
//PINGTIME PING timeout 
const PINGTIME = 5 * time.Second

// ENTER command will allow someone to enter in the network
func ENTER(conn net.Conn, ctx *utils.Context) error {
	if !ctx.IsMasterNode() {
		return errors.New("Can't let someone enter if I am not the master node")
	}

	remoteIP := utils.GetRemoteIP(conn)
	fmt.Println("Remote address entering network", remoteIP)
	ctx.AddNode(remoteIP)

	_, err := conn.Write([]byte(fmt.Sprintf("LEADER %s\n", ctx.Leader())))

	if err != nil {
		return err
	}

	eventlog.EventNewNode(remoteIP)

	return nil
}

// LEADER command will change leader
func LEADER(conn net.Conn, ctx *utils.Context, newLeader string) error {
	remoteIP := utils.GetRemoteIP(conn)
	if remoteIP != ctx.MasterNode() {
		return errors.New("Only master node can change the leader")
	}

	ctx.ChangeLeader(newLeader)
	eventlog.EventLeaderElected(newLeader)

	return nil
}

// SORT received a chunk, and decompress it sorting and sent it back to the master
func SORT(conn net.Conn, ctx *utils.Context, chunk string, id string) error {
	remoteIP := utils.GetRemoteIP(conn)
	if remoteIP != ctx.MasterNode() {
		return errors.New("Only master node can send an array for sorting")
	}

	chunkSlice, err := utils.UncompressChunk(chunk)
	if err != nil {
		return err
	}
	sort.Ints(chunkSlice)

	_, err = conn.Write([]byte(fmt.Sprintf("SORTED %s %s\n", utils.CompressChunk(chunkSlice), id)))

	if err != nil {
		return err
	}

	return nil
}

// WORK will receive an IP that is requesting work.
// If master will send an array for sorting
// If just leader will redirect to master
func WORK(conn net.Conn, ctx *utils.Context, remoteIP string) error {
	if ctx.IsMasterNode() {

		return nil
	}

	if ctx.IsLeader() {

		return nil
	}

	return errors.New("Only master node or leader can receive a WORK order")
}

// DIED reports that a node died
func DIED(IP string, ctx *utils.Context) (bool){ //Returns true if the dead node was the leader
	eventlog.EventDeadNode(IP)
	//ctx.RemoveNode(IP)
	return (ctx.Leader() == IP) 
}

// PING returns if nodes are connected
func PING(connection net.Conn, masterslave bool) (bool){ //True = master | False = slave
	if (masterslave) {
        timer := time.Now()
		for {
            fmt.Fprintf(connection, "PING")
			buffer := make([]byte, 1024)
            rcv,_ := connection.Read(buffer)
            if reflect.DeepEqual(buffer[:rcv], []byte("PONG")) {
                break
            } else if time.Since(timer) > PINGTIME {
                return false
            }
		}
	} else {
        timer := time.Now()
		for {
            buffer := make([]byte, 1024)
            rcv, _ := connection.Read(buffer)
            if reflect.DeepEqual(buffer[:rcv], []byte("PING")) {
                connection.Write([]byte("PONG"))
                break
            } else if time.Since(timer) > PINGTIME {
                return false
            }
		}
    }
    return true //Returns true if the conexion is OK
}

//HeartBeat periodically calls PING and return false if conexion fails
func HeartBeat(connection net.Conn, masterslave bool) (bool){
    for {
        timer := time.Now()
        for {
            if time.Since(timer) >= HEARTBTIME {
                if PING(connection, masterslave) {
                    break;
                } 
                return false
            }
        }
    }
}