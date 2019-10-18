// disable metamask

const web3 = new Web3(Web3.givenProvider || "ws://localhost:8546");
const shh = web3.shh;

const generatePayload = (text) => {
  const monthInMs = 60 * 60 * 24 * 31 * 1000;
  const timestamp = Math.round((new Date()).getTime());
  const payload = `["~#c4",["${text}","text/plain","~:public-group-user-message",${(timestamp + monthInMs) * 100},${timestamp}]]`

  return payload;
};

const chatName = "foo-bar-baz";
const chatNameHex = web3.utils.toHex(chatName);
const fullTopic = web3.utils.sha3(chatNameHex);
const topic = fullTopic.substr(0, 10);
const text = "test from a bot " + (new Date().getTime());
const payload = generatePayload(text);
const payloadHex = web3.utils.toHex(payload);
var symKeyID,
    keyPairID;

Promise.all([
  web3.shh.newKeyPair().then((id) => { keyPairID = id; }),
  web3.shh.generateSymKeyFromPassword(chatName).then((id) => { symKeyID = id; }),
  web3.shh.setMinPoW(0.002),
]).then(() => {
    web3.shh.subscribe("messages", {
      symKeyID: symKeyID,
      topics: [topic]
    }).on('data', function(data) {
      console.log("Received data: ", data)
      const message = web3.utils.toUtf8(data.payload)
      console.log(`Message received: ${message}`)
    });
}).then(() => {
  var msg = {
    symKeyID: symKeyID,
    sig: keyPairID,
    ttl: 10,
    topic: topic,
    payload: payloadHex,
    powTime: 1,
    powTarget: 0.002
  };

  console.log("Sending message", msg);

  web3.shh.post(msg)
    .then(h => console.log(`Message with hash ${h} successfuly sent`))
    .catch(err => console.log("Error: ", err));
});
