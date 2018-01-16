//
//Copyright [2016] [SnapRoute Inc]
//
//Licensed under the Apache License, Version 2.0 (the "License");
//you may not use this file except in compliance with the License.
//You may obtain a copy of the License at
//
//    http://www.apache.org/licenses/LICENSE-2.0
//
//       Unless required by applicable law or agreed to in writing, software
//       distributed under the License is distributed on an "AS IS" BASIS,
//       WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
//       See the License for the specific language governing permissions and
//       limitations under the License.
//
// _______  __       __________   ___      _______.____    __    ____  __  .___________.  ______  __    __
// |   ____||  |     |   ____\  \ /  /     /       |\   \  /  \  /   / |  | |           | /      ||  |  |  |
// |  |__   |  |     |  |__   \  V  /     |   (----` \   \/    \/   /  |  | `---|  |----`|  ,----'|  |__|  |
// |   __|  |  |     |   __|   >   <       \   \      \            /   |  |     |  |     |  |     |   __   |
// |  |     |  `----.|  |____ /  .  \  .----)   |      \    /\    /    |  |     |  |     |  `----.|  |  |  |
// |__|     |_______||_______/__/ \__\ |_______/        \__/  \__/     |__|     |__|      \______||__|  |__|
//

package server

import (
	"html/template"
	"infra/notifierd/alarmNotifier"
	"infra/notifierd/eventNotifier"
	"infra/notifierd/faultNotifier"
	"infra/notifierd/objects"
	"net/http"
	"utils/logging"
)

type NMGRServer struct {
	Logger        logging.LoggerIntf
	paramsDir     string
	ReqChan       chan *ServerRequest
	InitDone      chan bool
	EventEnable   bool
	FaultEnable   bool
	AlarmEnable   bool
	eventNotifier *eventNotifier.Notifier
	faultNotifier *faultNotifier.Notifier
	alarmNotifier *alarmNotifier.Notifier
}

func NewNMGRServer(initParams *ServerInitParams) *NMGRServer {
	nMgrServer := &NMGRServer{}
	nMgrServer.Logger = initParams.Logger
	nMgrServer.paramsDir = initParams.ParamsDir
	nMgrServer.InitDone = make(chan bool)
	nMgrServer.ReqChan = make(chan *ServerRequest)
	return nMgrServer
}

func (server *NMGRServer) InitServer() error {
	server.EventEnable = true
	server.FaultEnable = true
	server.AlarmEnable = true
	notifierPort, err := server.getNotifierPort()
	if err != nil {
		return err
	}
	dmnList, err := server.getDmnList()
	if err != nil {
		return err
	}
	notifierParams := &objects.NotifierParam{
		Logger:  server.Logger,
		DmnList: dmnList,
	}
	server.eventNotifier = eventNotifier.NewNotifier(notifierParams)
	server.faultNotifier = faultNotifier.NewNotifier(notifierParams)
	server.alarmNotifier = alarmNotifier.NewNotifier(notifierParams)
	go server.StartNotifier(notifierPort)
	return nil
}

func (server *NMGRServer) StartNotifier(port string) {
	http.HandleFunc("/events", server.eventNotifier.ProcessNotification)
	http.HandleFunc("/faults", server.faultNotifier.ProcessNotification)
	http.HandleFunc("/alarms", server.alarmNotifier.ProcessNotification)
	http.HandleFunc("/", server.loadHomePage)
	addr := ":" + port
	server.Logger.Err(http.ListenAndServe(addr, nil))
}

func (server *NMGRServer) StartServer() {
	err := server.InitServer()
	if err != nil {
		server.Logger.Err("Error Initializing server")
		return
	}
	server.InitDone <- true
	for {
		select {
		case req := <-server.ReqChan:
			server.Logger.Info("Server request received - ", *req)
			switch req.Op {
			case UPDATE_NOTIFIER_ENABLE:
				if val, ok := req.Data.(*UpdateNotifierEnableInArgs); ok {
					server.updateNotifierEnable(val.NotifierEnableOld, val.NotifierEnableNew, val.AttrSet)
				}
			default:
				server.Logger.Err("Error: Server received unrecognized request - ", req.Op)
			}
		}
	}
}

func (server *NMGRServer) loadHomePage(w http.ResponseWriter, r *http.Request) {
	homeTemplate.Execute(w, "")
}

var homeTemplate = template.Must(template.New("").Parse(`
<!DOCTYPE html>
<head>
<meta charset="utf-8">
<script>  
window.addEventListener("load", function(evt) {

    var eventsOutput = document.getElementById("eventsOutput");
    var faultsOutput = document.getElementById("faultsOutput");
    var alarmsOutput = document.getElementById("alarmsOutput");
    var eventsWs;
    var faultsWs;
    var alarmsWs;
    var host = window.location.host;
    var wsStr = "ws://";
    var baseAddr = wsStr.concat(host);

    var printEvent = function(message) {
        var d = document.createElement("div");
        var x = document.createElement("HR");
	var obj = JSON.parse(message)
        d.innerHTML = "OwnerName: " + obj.OwnerName + "<br>Description: " + obj.Description + "<br>OccuranceTime: " + obj.TimeStamp + "<br>Details: " + message;
        eventsOutput.insertBefore(x, eventsOutput.firstChild);
        eventsOutput.insertBefore(d, eventsOutput.firstChild);
    };

    var printFault = function(message) {
        var d = document.createElement("div");
        var x = document.createElement("HR");
	var obj = JSON.parse(message)
	if (obj.ResolutionTime == "N/A") {
        d.innerHTML = "OwnerName: " + obj.OwnerName + "<br>Description: " + obj.Description + "<br>Src Object:  {" + obj.SrcObjKey + "}<br>OccuranceTime: " + obj.OccuranceTime;
	} else {
        d.innerHTML = "Owner Name: " + obj.OwnerName + " <br>Description: " + obj.Description + "<br>Src Object:  {" + obj.SrcObjKey + "}<br>Occurance Time: " + obj.OccuranceTime + "<br>Resolution Time: " + obj.ResolutionTime + "<br>Resolution Reason: " + obj.ResolutionReason;
	}
        faultsOutput.insertBefore(x, faultsOutput.firstChild);
        faultsOutput.insertBefore(d, faultsOutput.firstChild);
    };

    var printAlarm = function(message) {
        var d = document.createElement("div");
        var x = document.createElement("HR");
	var obj = JSON.parse(message)
	if (obj.ResolutionTime == "N/A") {
        d.innerHTML = "OwnerName: " + obj.OwnerName + "<br>Description: " + obj.Description + "<br>Src Object:  {" + obj.SrcObjKey + "}<br>OccuranceTime: " + obj.OccuranceTime + "<br>Severity: " + obj.Severity;
	} else {
        d.innerHTML = "Owner Name: " + obj.OwnerName + "<br>Description: " + obj.Description + "<br>Src Object:  {" + obj.SrcObjKey + "}<br>Occurance Time: " + obj.OccuranceTime + "<br>Resolution Time: " + obj.ResolutionTime + "<br>Resolution Reason: " + obj.ResolutionReason + "<br>Severity: " + obj.Severity;
	}
        alarmsOutput.insertBefore(x, alarmsOutput.firstChild);
        alarmsOutput.insertBefore(d, alarmsOutput.firstChild);
    };

    document.getElementById("startEvents").onclick = function(evt) {
	var addr = baseAddr.concat("/events");
        if (eventsWs) {
            eventsWs.close();
            eventsWs = new WebSocket(addr);
            eventsWs.onopen = function(evt) {
            }
        } else {
            eventsWs = new WebSocket(addr);
            eventsWs.onopen = function(evt) {
            }
	}
        eventsWs.onclose = function(evt) {
            eventsWs = null;
        }
        eventsWs.onmessage = function(evt) {
            printEvent(evt.data);
        }
        eventsWs.onerror = function(evt) {
            printEvent("ERROR: " + evt.data);
        }
        return false;
    };

    document.getElementById("startFaults").onclick = function(evt) {
	var addr = baseAddr.concat("/faults");
        if (faultsWs) {
            faultsWs.close();
            faultsWs = new WebSocket(addr);
            faultsWs.onopen = function(evt) {
            }
        } else {
            faultsWs = new WebSocket(addr);
            faultsWs.onopen = function(evt) {
            }
	}
        faultsWs.onclose = function(evt) {
            faultsWs = null;
        }
        faultsWs.onmessage = function(evt) {
            printFault(evt.data);
        }
        faultsWs.onerror = function(evt) {
            printFault("ERROR: " + evt.data);
        }
        return false;
    };

    document.getElementById("startAlarms").onclick = function(evt) {
	var addr = baseAddr.concat("/alarms");
        if (alarmsWs) {
            alarmsWs.close();
            alarmsWs = new WebSocket(addr);
            alarmsWs.onopen = function(evt) {
            }
        } else {
            alarmsWs = new WebSocket(addr);
            alarmsWs.onopen = function(evt) {
            }
	}
        alarmsWs.onclose = function(evt) {
            alarmsWs = null;
        }
        alarmsWs.onmessage = function(evt) {
            printAlarm(evt.data);
        }
        alarmsWs.onerror = function(evt) {
            printAlarm("ERROR: " + evt.data);
        }
        return false;
    };

    document.getElementById("stopEvents").onclick = function(evt) {
        if (!eventsWs) {
            return false;
        }
        eventsWs.close();
        return false;
    };

    document.getElementById("stopFaults").onclick = function(evt) {
        if (!faultsWs) {
            return false;
        }
        faultsWs.close();
        return false;
    };

    document.getElementById("stopAlarms").onclick = function(evt) {
        if (!alarmsWs) {
            return false;
        }
        alarmsWs.close();
        return false;
    };

});
</script>
<style type="text/css">
header, footer {
    padding: 1px;
    color: white;
    background-color: black;
    clear: left;
    text-align: center;
}
ul.tab {
    list-style-type: none;
    margin: 0;
    padding: 0;
    overflow: hidden;
    border: 1px solid #ccc;
    background-color: #f1f1f1;
}

ul.tab li {float: left;}

ul.tab li a {
    display: inline-block;
    color: black;
    text-align: center;
    padding: 14px 16px;
    text-decoration: none;
    transition: 0.3s;
    font-size: 17px;
}

ul.tab li a:hover {background-color: #ddd;}

ul.tab li a:focus, .active {background-color: #ccc;}

div.scroll {
    background-color: #d6d6c2;
    height: 370px;
    width: 100%;
    overflow-y: auto;
}

table
{
    table-layout: fixed;
    width: 100%;
}

.heading {
    text-align: center;
    background-color: #e6ffe6;
}

</style>
</head>
<body>
<header>
   <h1>SnapRoute's Flexswitch Events, Faults and Alarms</h1>
</header>

<ul class="tab">
  <li><a href="#" class="tablinks" id="startAlarms">Start Monitoring Alarms</a></li>
  <li><a href="#" class="tablinks" id="startFaults">Start Monitoring Faults</a></li>
  <li><a href="#" class="tablinks" id="startEvents">Start Monitoring Events</a></li>
  <li><a href="#" class="tablinks" id="stopAlarms">Stop Monitoring Alarms</a></li>
  <li><a href="#" class="tablinks" id="stopFaults">Stop Monitoring Faults</a></li>
  <li><a href="#" class="tablinks" id="stopEvents">Stop Monitoring Events</a></li>
</ul>

<p class="heading"> <font size="5"> Alarms </font> <p>
<div class="scroll">
<table>
<tr>
<td>
<div id="alarmsOutput">
</div>
</td>
</tr>
</table>
</div>

<p class="heading"> <font size="5"> Faults </font> <p>
<div class="scroll">
<table>
<tr>
<td>
<div id="faultsOutput">
</div>
</td>
</tr>
</table>
</div>

<p class="heading"> <font size="5"> Events </font> </p>
<div class="scroll">
<table>
<tr>
<td>
<div id="eventsOutput">
</div>
</td>
</tr>
</table>
</div>

<footer>Copyright © snaproute.com</footer>
</body>
</html>
`))
