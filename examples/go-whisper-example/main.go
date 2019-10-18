package main

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/ethereum/go-ethereum/whisper/shhclient"
	"github.com/ethereum/go-ethereum/whisper/whisperv6"
	"golang.org/x/crypto/sha3"
)

func check(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func makePayload(text string) string {
	timestamp := time.Now().Unix() * 1000
	format := `["~#c4",["%s","text/plain","~:public-group-user-message",%d,%d]]`
	oneMonthInMs := int64(60 * 60 * 24 * 31 * 1000)
	payload := fmt.Sprintf(format, text, (timestamp+oneMonthInMs)*100, timestamp)

	return payload
}

func topicFromChatName(chatName string) whisperv6.TopicType {
	h := sha3.NewLegacyKeccak256()
	h.Write([]byte(chatName))
	fullTopic := h.Sum(nil)

	return whisperv6.BytesToTopic(fullTopic)
}

func readLine() string {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print(">> ")
	text, err := reader.ReadString('\n')
	check(err)

	return strings.TrimSpace(text)
}

func main() {
	ipcPath := os.Getenv("IPC_PATH")
	if ipcPath == "" {
		fmt.Printf("you must specify the IPC_PATH environment variable. ")
		fmt.Printf("Try with: export IPC_PATH=path-to-file.ipc\n")
		os.Exit(1)
	}

	if len(os.Args) != 2 {
		fmt.Printf("you must specify the public chat you want to join.\n")
		fmt.Printf("Try with: %s CHAT_NAME\n", os.Args[0])
		os.Exit(1)
	}

	// public chat
	chatName := os.Args[1]

	whisperClient, err := shhclient.Dial(ipcPath)
	check(err)

	ctx := context.Background()

	err = whisperClient.SetMinimumPoW(ctx, 0.002)
	if err != nil {
		log.Fatal(err)
	}

	// topic
	topic := topicFromChatName(chatName)

	// symkey
	symKeyID, err := whisperClient.GenerateSymmetricKeyFromPassword(ctx, chatName)
	check(err)

	// keypair
	keyPairID, err := whisperClient.NewKeyPair(ctx)
	check(err)

	crit := whisperv6.Criteria{
		SymKeyID: symKeyID,
		Topics:   []whisperv6.TopicType{topic},
	}
	ch := make(chan *whisperv6.Message, 0)
	whisperClient.SubscribeMessages(ctx, crit, ch)
	go func() {
		for {
			m := <-ch
			fmt.Printf("\nRECEIVED: [%x] \n%s\n\n>>", m.Sig, m.Payload)
		}
	}()

	for {
		text := readLine()

		// payload
		payload := makePayload(text)

		msg := whisperv6.NewMessage{
			SymKeyID:  symKeyID,
			Sig:       keyPairID,
			TTL:       15,
			Topic:     topic,
			Payload:   []byte(payload),
			PowTime:   5,
			PowTarget: 0.002,
		}

		fmt.Printf("MESSAGE: %+v\n", msg)

		hash, err := whisperClient.Post(ctx, msg)
		if err != nil {
			log.Fatalf("ERROR: %+v", err)
		}
		fmt.Printf("-- sent message with has %s\n", hash)
	}
}
