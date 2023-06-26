import { check,sleep } from 'k6';
import http from 'k6/http';

let shortUrl;
export function setup() {
    var url ="http://127.0.0.1/api/v1/urls"
    var payload = JSON.stringify({
        "url": "https://www.google.com",
        "expireAt": "2024-09-09T00:00:00Z"
    });
    var headers = {
        'Content-Type': 'application/json' ,
    }
    const response = http.post(url, payload,{ headers: headers });
    shortUrl = response.json().shortUrl
    return {
        shortUrl,
    };
}

export default function (data) {
    const { shortUrl } = data;

    const resp = http.get(shortUrl);
    if (resp.status!==200){
        console.log(resp)
    }
    check(resp, {
        "status is 200": (res) => res.status === 200,
      });
}
