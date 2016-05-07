import * as http from 'http';
import * as url from 'url';
import * as _ from 'lodash';
import * as qs from 'querystring';
import * as momentOrg from 'moment';
const moment: () => IMoment = require('moment-timezone');


import RequestOptions = http.RequestOptions;
import ServerResponse = http.ServerResponse;
import IncomingMessage = http.IncomingMessage;


interface IMoment extends momentOrg.Moment {
  tz(ofset: string|number): IMoment;
}

/**
 * 1: recive post request to process payment
 * 2: (Time pases and the user inputs data to the system)
 * 3: System responds with sample data from paypal (may be customized later)
 * 4: Verify empty response and status code 200 (DONE)
 * 5: Wait for reciving message from client
 * 6: Verify that message contains all fields of the original message and the extra "cmd=_notify-validate" as the first field (preceded field)
 * 7: Send "VERIFIED" / "INVALID" message
 */

const configuration = {
  serverPort: 8081,
  serverProcessPaymentUrl: '/processPayment',
  serverIpnCallback: '/cgi-bin/webscr',

  //Client Settings  
  ipnAddress: 'localhost',
  ipnPort: 8080,
  ipnPath: '/rest/paypal/ipn'
}

let processIpnResponse: (req: IncomingMessage, res: ServerResponse) => void = null;

http.createServer((req, res) => {

  if (req.url === configuration.serverProcessPaymentUrl) {
    console.log('')
    sendIpnDataMessage(req, res).then((sentMessage: string) => {
      res.write('<a href="//localhost:9000/status">back</a>\n');

      log('Data message sent', res);
      log('server closed connection', res);
      log('status code 200 and empty response (OK)', res);
      processIpnResponse = waitIpnDataMessageResponse(sentMessage, res);
    }, (errMsg: string) => {
      log(errMsg, res);
      res.end();
    });
  } else if (req.url === configuration.serverIpnCallback) {
    if (processIpnResponse === null) {
      console.log('recived ipn response, while not expecting one');
      res.end();
    } else {
      processIpnResponse(req, res);
      processIpnResponse = null;
    }
  } else {
    if (req.url !== '/favicon.ico') {
      console.log('Ignoring request: ' + req.url);
    }
    res.end();
  }

}).listen(configuration.serverPort, () => {
  console.log("listening for requests on port: " + configuration.serverPort);
});

//Reset on timeout and print error
//READ POST DATA must be "cmd=_notify-validate" + waitingForMessage
//Send "VERIFIED" / "INVALID" message 

function waitIpnDataMessageResponse(expectedMsg: string, res: ServerResponse): (req: IncomingMessage, res: ServerResponse) => void {

  let isPending = true;

  setTimeout(() => {
    if (isPending) {
      isPending = false;
      log("Timeout: Message was not recived in time", res);
      res.end();
    }
    processIpnResponse = null;
  }, 5000);

  return function messageCB(req: IncomingMessage, newRes: ServerResponse) {
    if (!isPending) {
      log("Message late: Message was not recived in time");
      isPending = false;
      return;
    }
    isPending = false;
    const extraValues = "cmd=_notify-validate&";

    readAllData(req).then((recivedMsg) => {
      const messageValid = (recivedMsg === extraValues + expectedMsg);

      if (res.statusCode !== 200) {
        log('Status code was not 200 it was: ' + res.statusCode + '\n' + recivedMsg);
      } else if (messageValid) {
        log("Success VERIFIED", res);
        newRes.write("VERIFIED");
      } else {
        log("failure INVALID", res);
        log("Recived: " + recivedMsg, res);
        log("Expected: " + extraValues + expectedMsg, res);
        newRes.write("INVALID");
      }

      log('<br>', res)

      recivedMsg.split('&').forEach((line) => {
        log(`${line}&`, res);
      });

      newRes.end();
      res.end();
    });

  }
}

function sendIpnDataMessage(clientReq: IncomingMessage, serverRes: ServerResponse) {

  return new Promise<string>((resole, reject) => {

    readAllData(clientReq).then((body) => {
      const dynamicFields = {
        "payment_date": moment().tz("Europe/Copenhagen").utcOffset('+08:00').format('HH:mm:ss MMM D, YYYY z') //"20:12:59 Jan 13, 2009 PST",
      };
      const formQueryObj = qs.parse(body)
      const queryObj = url.parse(clientReq.url, true).query;

      const ipnBody = [paypalSampleIpnMsg, queryObj, formQueryObj, dynamicFields].reduce((prev, cur) => {
        return _.merge(prev, cur);
      }, {});

      const ipnBodyEncoded = encodePostData(ipnBody);

      const options: RequestOptions = {
        hostname: configuration.ipnAddress,
        port: configuration.ipnPort,
        path: configuration.ipnPath,
        method: 'POST',
        headers: {
          'Content-Type': 'application/x-www-form-urlencoded',
          'Content-Length': ipnBodyEncoded.length
        }
      }

      // https://developer.paypal.com/webapps/developer/docs/classic/ipn/integration-guide/IPNIntro/#id08CKFJ00JYK
      const req = http.request(options, (res: IncomingMessage) => {
        readAllData(res).then((responseBody) => {
          if (res.statusCode !== 200) {
            reject('Status code was not 200 it was: ' + res.statusCode + '\n' + responseBody);
          } else if (responseBody != '') {
            reject('Body was not empty it contained: ' + responseBody)
          } else {
            resole(ipnBodyEncoded);
          }
        });
      });

      req.on('error', (err: Error) => {
        reject(err.message);
      });

      req.write(ipnBodyEncoded);

      req.end();
    });
  });
}









const paypalSampleIpnMsg = {
  "mc_gross": "19.95",
  "protection_eligibility": "Eligible",
  "address_status": "confirmed",
  "payer_id": "LPLWNMTBWMFAY",
  "tax": "0.00",
  "address_street": "1 Main St",
  "payment_date": "20:12:59 Jan 13, 2009 PST",
  "payment_status": "Completed",
  "charset": "windows-1252",
  "address_zip": "95131",
  "first_name": "Test",
  "mc_fee": "0.88",
  "address_country_code": "US",
  "address_name": "Test User",
  "notify_version": "2.6",
  "custom": "",
  "payer_status": "verified",
  "address_country": "United States",
  "address_city": "San Jose",
  "quantity": "1",
  "verify_sign": "AtkOfCXbDm2hu0ZELryHFjY-Vb7PAUvS6nMXgysbElEn9v-1XcmSoGtf",
  "payer_email": "gpmac_1231902590_per@paypal.com",
  "txn_id": "61E67681CH3238416",
  "payment_type": "instant",
  "last_name": "User",
  "address_state": "CA",
  "receiver_email": "gpmac_1231902686_biz@paypal.com",
  "payment_fee": "0.88",
  "receiver_id": "S8XGHLYDW9T3S",
  "txn_type": "express_checkout",
  "item_name": "",
  "mc_currency": "DKK",
  "item_number": "",
  "residence_country": "DK",
  "test_ipn": "1",
  "handling_amount": "0.00",
  "transaction_subject": "",
  "payment_gross": "19.95",
  "shipping": "0.00"
}





function testParsing() {
  const bodySampleMessage = 'mc_gross=19.95&protection_eligibility=Eligible&address_status=confirmed&payer_id=LPLWNMTBWMFAY&tax=0.00&address_street=1+Main+St&payment_date=20%3A12%3A59+Jan+13%2C+2009+PST&payment_status=Completed&charset=windows-1252&address_zip=95131&first_name=Test&mc_fee=0.88&address_country_code=US&address_name=Test+User&notify_version=2.6&custom=&payer_status=verified&address_country=United+States&address_city=San+Jose&quantity=1&verify_sign=AtkOfCXbDm2hu0ZELryHFjY-Vb7PAUvS6nMXgysbElEn9v-1XcmSoGtf&payer_email=gpmac_1231902590_per%40paypal.com&txn_id=61E67681CH3238416&payment_type=instant&last_name=User&address_state=CA&receiver_email=gpmac_1231902686_biz%40paypal.com&payment_fee=0.88&receiver_id=S8XGHLYDW9T3S&txn_type=express_checkout&item_name=&mc_currency=USD&item_number=&residence_country=US&test_ipn=1&handling_amount=0.00&transaction_subject=&payment_gross=19.95&shipping=0.00';

  const body = parsePostData(bodySampleMessage);
  const encodedBody = encodePostData(body);

  log(`Are equal: ${bodySampleMessage === encodedBody}`);
}

function parsePostData(data: string) {
  return data.split('&').map(pair => pair.split('='))
    .map(pair => pair.map((part) => decodeURIComponent(part.replace(/\+/g, ' '))))
    .reduce(
    (prev: any, next: Array<string>): [{ string: string }] => {
      prev[next[0]] = next[1];
      return prev;
    }, {});
}

function encodePostData(data: any) {

  return Object.keys(data)
    .map((key: string) => [key, data[key]])
    .map((pair: Array<string>) => pair.map(part => encodeURIComponent(part).replace(/\%20/g, '+')))
    .map((pair: Array<string>) => pair.join('='))
    .join('&');
}

function log(msg: string, res?: ServerResponse) {

  if (res && res.writable) {
    res.write(`<div>${msg}</div>\n`);
  }
  console.log(msg);

}


function readAllData(res: IncomingMessage): Promise<string> {
  return new Promise((resolve, reject) => {
    let fullBody = '';

    res.on('data', (chunk: Buffer) => {
      fullBody += chunk.toString();
    });

    res.on('end', () => {
      resolve(fullBody);
    });
  });
}