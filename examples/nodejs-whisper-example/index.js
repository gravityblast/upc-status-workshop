const StatusJS = require('status-js-api');

const readline = require('readline').createInterface({
  input: process.stdin,
  output: process.stdout
})

if (process.argv.length != 3) {
  console.log("node index.js CHAT_NAME");
  process.exit(1);
}

let chatName,
    status;

const send = (text) => {
  status.sendMessage(`#${chatName}`, text);
}

const prompt = () => {
  readline.question(">> ", (text) => {
    send(text)
  })
}

const main = async () => {
  status = new StatusJS();
  chatName = process.argv[2];

  await status.connect("http://localhost:8545");
  await status.joinChat(chatName);

  status.onMessage(chatName, (err, data) => {
    readline.pause();
    console.log()
    if(err) {
      console.error("Error: " + err + "\n");
    } else {
      console.log(`message received from ${data.username}:`);
      console.log(data.payload);
    }

    readline.prompt()
    readline.resume()
    prompt(true);
  });

  prompt(false);
}

main();
