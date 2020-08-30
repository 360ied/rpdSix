package bigemojicommand

import (
	"bytes"
	"encoding/base64"
	"errors"
	"fmt"
	"image"
	_ "image/gif"
	_ "image/jpeg"
	"image/png"
	"rpdSix/commands"
	"rpdSix/helpers"
	"strconv"
)

func Initialize() {
	commands.AddCommand(
		commands.Command{
			Run:                         run,
			Names:                       []string{"bigemoji"},
			ExpectedPositionalArguments: []string{emojiNameArg, gridSizeArg},
		},
	)
}

const (
	emojiNameArg    = "emojiNameArg"
	gridSizeArg     = "gridSizeArg"
	defaultGridSize = 2
)

func run(ctx commands.CommandContext) error {
	if len(ctx.Message.Attachments) == 0 {
		var _, err = ctx.Session.ChannelMessageSend(ctx.Message.ChannelID, "No attachments were found!")
		return err
	}

	var attachment = ctx.Message.Attachments[0]
	// consider using ProxyURL instead of URL
	resp, err := ctx.Session.Client.Get(attachment.URL)
	if err != nil {
		return err
	} else if resp.StatusCode != 200 {
		return errors.New(fmt.Sprint("status code is not 200, status code is: ", resp.StatusCode))
	}

	var img, _, err2 = image.Decode(resp.Body)
	if err2 != nil {
		return err
	}

	var gridSize int

	var gridSizeStr, exists = ctx.Arguments[gridSizeArg]

	if !exists {
		gridSize = defaultGridSize
	} else {
		var gridSizeTmp, err3 = strconv.Atoi(gridSizeStr)
		if err3 != nil {
			return err3
		}
		gridSize = gridSizeTmp
	}

	var grid [][]image.Image

	var imgSize = img.Bounds().Size()

	var partitionSizeX = imgSize.X / gridSize
	var partitionSizeY = imgSize.Y / gridSize
	//fmt.Println(partitionSizeY)
	//fmt.Println(partitionSizeY)

	for i := 0; i < gridSize; i++ {
		var row []image.Image
		for j := 0; j < gridSize; j++ {
			row = append(row, img.(helpers.SubImager).SubImage(
				image.Rect(
					partitionSizeX*j, partitionSizeY*i,
					partitionSizeX*j+partitionSizeX, partitionSizeY*i+partitionSizeY)))
			//fmt.Println(fmt.Sprint(partitionSizeX*j, partitionSizeY*i,
			//	partitionSizeX*j, partitionSizeY*i))
		}
		grid = append(grid, row)
	}

	// make emojis

	//var emojis [][]*discordgo.Emoji

	var emojiBaseName, exists2 = ctx.Arguments[emojiNameArg]
	if !exists2 {
		var _, err3 = ctx.Session.ChannelMessageSend(ctx.Message.ChannelID, "Emoji name not found!")
		return err3
	}

	//fmt.Println(emojiBaseName)
	//fmt.Println(len(emojiBaseName))

	var messageString = ""

	for i := 0; i < gridSize; i++ {
		//var row []*discordgo.Emoji
		for j := 0; j < gridSize; j++ {
			//fmt.Println(i, j)

			var buffer bytes.Buffer

			//goland:noinspection GoNilness
			var err4 = png.Encode(&buffer, grid[i][j])
			if err4 != nil {
				return err4
			}

			//fmt.Println("encoding image")

			var encodedImage = base64.StdEncoding.EncodeToString(buffer.Bytes())

			//fmt.Println(encodedImage)
			//fmt.Println("length of encoded image", len(encodedImage))
			//fmt.Println(fmt.Sprint(emojiBaseName, "_", i, "_", j))

			var emoji, err5 = ctx.Session.GuildEmojiCreate(
				ctx.Message.GuildID,
				// humans naturally start counting at 1, not 0
				fmt.Sprint(emojiBaseName, "_", i+1, "_", j+1),
				fmt.Sprint("data:png;base64,", encodedImage),
				nil)
			if err5 != nil {
				return err5
			}

			//fmt.Println(emoji.Name)

			messageString += fmt.Sprint(":", emoji.Name, ":")

			//row = append(row, emoji)
		}
		messageString += "\n"
		//emojis = append(emojis, row)
	}

	var _, err3 = ctx.Session.ChannelMessageSend(ctx.Message.ChannelID, messageString)

	return err3
}
