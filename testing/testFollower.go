package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
)

type Follower struct {
	CurrentTerm int
	Logs        []Log
	Timeout     int
}

func (this *Follower) Init() {
	this.CurrentTerm = 4
	this.Logs = make([]Log, 0)
	this.Logs = append(this.Logs, Log{1, 1, "add"})
	this.Logs = append(this.Logs, Log{1, 2, "add"})
	this.Logs = append(this.Logs, Log{1, 3, "add"})
	this.Logs = append(this.Logs, Log{4, 4, "add"})
	this.Logs = append(this.Logs, Log{4, 5, "add"})
	this.Logs = append(this.Logs, Log{5, 6, "add"})
	this.Logs = append(this.Logs, Log{5, 7, "add"})
	this.Logs = append(this.Logs, Log{6, 8, "add"})
	this.Logs = append(this.Logs, Log{6, 9, "add"})
	this.Logs = append(this.Logs, Log{6, 10, "add"})
	this.Logs = append(this.Logs, Log{7, 11, "add"})
	this.Logs = append(this.Logs, Log{7, 12, "add"})

	this.Timeout = 100
}

func (this *Follower) CheckPrev(index, term int) bool {
	if index > len(this.Logs) {
		return false
	} else if index == 0 {
		return len(this.Logs) == 0
	} else if len(this.Logs) == 0 {
		return true
	}

	log := this.Logs[index-1]
	return log.Term == term
	// if log.Term != term {
	// 	return false
	// }

	// return true
}

func (this *Follower) PrintLogs() {
	for i := 0; i < len(this.Logs); i++ {
		fmt.Println(this.Logs[i].ToStr())
	}
}

func (this *Follower) HandleAppendEntriesRPC(requests *AppendReqs, responses *Responses) {
	for _, rpc := range requests.Rpcs {

		// validity check ?

		if rpc.Term < this.CurrentTerm {
			responses.Resps = append(responses.Resps, AppendResp{this.CurrentTerm, false})
			continue
		}

		//if rpc.Term > this.CurrentTerm {
		//	this.CurrentTerm = rpc.Term
		//}

		if len(rpc.Entries) == 0 {
			// TODO: reset timer
			continue
		}

		// return failure if log does not contain an entry at prevLogIndex whose
		//  term matches prevLogTerm
		if !this.CheckPrev(rpc.PrevLogIndex, rpc.PrevLogTerm) {
			responses.Resps = append(responses.Resps, AppendResp{this.CurrentTerm, false})
			continue
		}

		if rpc.Term > this.CurrentTerm {
			this.CurrentTerm = rpc.Term
		}

		// check for conflicts
		// need to sort the entries first
		// sort.SliceStable(rpc.Entries, func(i, j int) bool {
		// 	log_i := rpc.Entries[i]
		// 	log_j := rpc.Entries[j]

		// 	if log_i.Term < log_j.Term {
		// 		return true
		// 	} else if log_i.Term == log_j.Term {
		// 		return log_i.Index < log_j.Index
		// 	}

		// 	return false
		// })

		// conflict_index := -1
		// for i := rpc.PrevLogIndex; i < len(this.Logs); i++ {
		// 	if this.Logs[i].Term != rpc.Entries[0].Term || this.Logs[i].Index != rpc.Entries[0].Index || strings.Compare(rpc.Entries[0].Command, this.Logs[i].Command) != 0 {
		// 		conflict_index = i
		// 		break
		// 	}
		// }

		// if conflict_index != -1 {
		// 	this.Logs = this.Logs[:conflict_index]
		// }

		// add all entries from RPC to follower's log
		// deepcopy them to prevent some issues
		// fmt.Println("Add new entry")
		if rpc.PrevLogIndex < len(this.Logs)-1 {
			this.Logs = this.Logs[:rpc.PrevLogIndex]
		}
		this.Logs = append(this.Logs, DeepCopyLogs(rpc.Entries)...)
		responses.Resps = append(responses.Resps, AppendResp{this.CurrentTerm, true})
	}
}

func main() {
	//data, err := ParseAppendReqFromFile("test_input/test_follower_same.json")
	data, err := ParseRPCAndRespFromFile("test_input/test_extra.json")
	reqs := data.Rpcs
	//resp := data.Resps
	if err == nil {
		//PrintAppendReqs(&reqs)
	} else {
		fmt.Println("err")
	}

	responses := Responses{}
	follower := Follower{}
	follower.Init()

	req_arr := AppendReqs{}
	req_arr.Rpcs = reqs
	follower.HandleAppendEntriesRPC(&req_arr, &responses)

	// open an output file
	out_file, out_file_err := os.OpenFile("follower_out_extra.json", os.O_CREATE|os.O_WRONLY, 0777)
	if out_file_err != nil {
		panic(out_file_err)
	}

	PrintResps(&responses)
	follower.PrintLogs()
	//PrintResps(&resp)

	writter := bufio.NewWriter(out_file)
	barr, json_err := json.MarshalIndent(responses, "", "    ")
	// fmt.Print(barr)
	// fmt.Println(len(responses.Resps))
	if json_err != nil {
		panic(json_err)
	}
	if _, write_err := writter.Write(barr); err != nil {
		panic(write_err)
	}
	writter.Flush()
}
