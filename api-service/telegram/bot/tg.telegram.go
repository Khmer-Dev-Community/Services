package bot

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"sync"
	"time" // Added for crypto/rand

	// Added for encoding/binary
	"crypto/rand"     // Added for crypto/rand
	"encoding/binary" // Added for encoding/binary

	"github.com/gotd/td/session"
	"github.com/gotd/td/telegram"
	"github.com/gotd/td/telegram/auth"
	"github.com/gotd/td/tg"

	// Your existing imports

	"telegram-service/telegram/services"
)

// NOTE: MessageData struct is not used in the provided logic, but kept if you need it elsewhere.
type MessageData struct {
	ChatID          string
	ChatTitle       string
	ChatType        string
	SenderID        string
	SenderUsername  string
	SenderFirstName string
	SenderLastName  string
	MessageText     string
	Timestamp       time.Time
	ReceivedAt      time.Time
	IsBot           bool
	RawUpdate       string
}

// BotAccount represents the userbot client.
type BotAccount struct {
	Name               string
	Debug              bool
	AppID              int
	APIHash            string
	PhoneNumber        string
	SessionFile        string
	mu                 sync.Mutex
	TransactionService services.TransactionService
	Client             *telegram.Client
	Dispatcher         *tg.UpdateDispatcher
}

// NewTgService is the constructor for the BotAccount (userbot).
func NewTgService(name string, appID int, apiHash, phoneNumber, sessionFile string, debug bool, isInteractive bool, trxService services.TransactionService) (*BotAccount, error) {
	// Create a new dispatcher
	dispatcher := tg.NewUpdateDispatcher()
	b := &BotAccount{
		Name:               name,
		Debug:              debug,
		AppID:              appID,
		APIHash:            apiHash,
		PhoneNumber:        phoneNumber,
		SessionFile:        sessionFile,
		TransactionService: trxService,
		Dispatcher:         &dispatcher,
	}
	//dispatcher.OnNewChannelMessage(b.onNewChannelMessage)
	//dispatcher.OnNewMessage(b.onNewMessage)
	client := telegram.NewClient(appID, apiHash, telegram.Options{
		SessionStorage: &session.FileStorage{Path: sessionFile},
		UpdateHandler:  dispatcher,
		// Logger:         log.New(os.Stderr, "gotd: ", log.LstdFlags), // Uncomment for verbose logging
	})

	b.Client = client
	return b, nil
}

// Start initiates the userbot connection and message processing.
func (b *BotAccount) Start() error {
	log.Printf("Starting userbot %s (gotd/td)...", b.Name)

	if err := b.Client.Run(context.Background(), func(ctx context.Context) error {
		status, err := b.Client.Auth().Status(ctx)
		if err != nil {
			return fmt.Errorf("auth status failed: %w", err)
		}
		if !status.Authorized {
			log.Println("Userbot not authenticated. Starting authentication flow...")
			if err := b.authenticate(ctx); err != nil {
				return fmt.Errorf("userbot authentication failed: %w", err)
			}
			log.Println("Userbot authentication successful!")
		} else {
			log.Printf("Userbot already authenticated as %s.", status.User.Username)
		}

		log.Printf("Userbot %s is now listening for updates...", b.Name)
		<-ctx.Done()
		return ctx.Err()
	}); err != nil {
		log.Fatalf("Userbot client run failed: %v", err)
	}
	return nil
}

func (b *BotAccount) authenticate(ctx context.Context) error {
	codeAuthenticator := auth.CodeAuthenticatorFunc(func(ctx context.Context, sentCode *tg.AuthSentCode) (string, error) {
		fmt.Printf("Enter the authentication code for %s: ", b.PhoneNumber)
		code, err := bufio.NewReader(os.Stdin).ReadString('\n')
		if err != nil {
			return "", err
		}
		return strings.TrimSpace(code), nil
	})
	flow := auth.NewFlow(
		auth.Constant(b.PhoneNumber, "", codeAuthenticator),
		auth.SendCodeOptions{},
	)
	return b.Client.Auth().IfNecessary(ctx, flow)
}
func (b *BotAccount) SendMessageToGroup(ctx context.Context, chatID int64, message string) error {
	peer, err := b.resolvePeer(ctx, chatID)
	if err != nil {
		return fmt.Errorf("failed to resolve peer for chat ID %d: %w", chatID, err)
	}
	var randomID int64
	if err := binary.Read(rand.Reader, binary.LittleEndian, &randomID); err != nil {
		return fmt.Errorf("failed to generate random ID for message: %w", err)
	}
	_, err = b.Client.API().MessagesSendMessage(ctx, &tg.MessagesSendMessageRequest{
		Peer:     peer,
		Message:  message,
		RandomID: randomID,
	})
	if err != nil {
		return fmt.Errorf("failed to send message to chat ID %d: %w", chatID, err)
	}

	log.Printf("Message sent to chat ID %d: %s", chatID, message)
	return nil
}

func (b *BotAccount) resolvePeer(ctx context.Context, chatID int64) (tg.InputPeerClass, error) {
	if chatID > 0 {

		users, err := b.Client.API().UsersGetUsers(ctx, []tg.InputUserClass{&tg.InputUser{UserID: chatID}})
		if err != nil || len(users) == 0 {
			return nil, fmt.Errorf("could not get user %d: %w", chatID, err)
		}
		user, ok := users[0].(*tg.User) // Type assertion: check if it's a *tg.User
		if !ok {
			return nil, fmt.Errorf("resolved user %d is not a concrete user (type %T)", chatID, users[0])
		}
		return &tg.InputPeerUser{
			UserID:     user.ID,
			AccessHash: user.AccessHash,
		}, nil
	} else if chatID < 0 && chatID > -1000000000000 { // Basic group chat ID starts with - (e.g., -123456789)

		return &tg.InputPeerChat{ChatID: -chatID}, nil
	} else if chatID < -1000000000000 { // Supergroup/channel ID starts with -100 (e.g., -100123456789)
		channelID := -chatID // Make it positive
		strChannelID := strconv.FormatInt(channelID, 10)
		if len(strChannelID) > 3 && strChannelID[0:3] == "100" {
			parsedID, err := strconv.ParseInt(strChannelID[3:], 10, 64)
			if err == nil {
				channelID = parsedID
			}
		}
		channelsResult, err := b.Client.API().ChannelsGetChannels(ctx, []tg.InputChannelClass{&tg.InputChannel{ChannelID: channelID}})
		if err != nil {
			return nil, fmt.Errorf("ChannelsGetChannels API call failed for channel %d: %w", channelID, err)
		}

		var chats []tg.ChatClass // This will hold the actual slice of chats
		switch v := channelsResult.(type) {
		case *tg.MessagesChats:
			chats = v.Chats
		case *tg.MessagesChatsSlice:
			chats = v.Chats
		default:
			return nil, fmt.Errorf("unexpected response type for ChannelsGetChannels: %T", channelsResult)
		}

		if len(chats) == 0 {
			return nil, fmt.Errorf("could not find channel %d in response", channelID)
		}

		// Now, we can safely access chats[0] and perform the type assertion to *tg.Channel
		channel, ok := chats[0].(*tg.Channel)
		if !ok {
			return nil, fmt.Errorf("resolved peer for channel %d is not a channel (type %T)", channelID, chats[0])
		}
		// --- END CORRECTED PART ---

		return &tg.InputPeerChannel{
			ChannelID:  channel.ID,
			AccessHash: channel.AccessHash,
		}, nil
	}
	return nil, fmt.Errorf("invalid chat ID format: %d", chatID)
}
