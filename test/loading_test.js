import { check,sleep } from 'k6';
import http from 'k6/http';


  

export default function () {
    var url ="http://127.0.0.1/api/v1/urls"
    var payload = JSON.stringify({
        "url": "https://www.google.com",
        "expireAt": "2024-09-09T00:00:00Z"
    });
    var headers = {
        'Content-Type': 'application/json' ,
    }
    const response = http.post(url, payload,{ headers: headers });
    if (response.status!==200){
        console.log(response.json())
    }
    check(response, {
        "status is 200": (res) => res.status === 200,
        "check body": (res) => res.json().id !== "",
      });
}
