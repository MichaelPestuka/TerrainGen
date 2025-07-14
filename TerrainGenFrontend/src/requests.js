import { startRenderer } from "./renderer";

var req = 1

export function getNewTerrain()
{
    if(req > 1)
    {
        return
    }
    req++;
    let new_req = new XMLHttpRequest();
    new_req.addEventListener("load", request_callback);
    new_req.open("QUERY", "http://localhost:8080", true);
    // new_req.setRequestHeader("Access-Control-Allow-Methods", "*");
    new_req.setRequestHeader("Content-Type", "application/json;charset=UTF-8");
    new_req.send(JSON.stringify({"width": 81, "height": 81}));
}

function request_callback()
{
    var parsed = JSON.parse(this.responseText)
    console.log(parsed)
    startRenderer(true, parsed)
}