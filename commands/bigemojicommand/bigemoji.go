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

	for i := 0; i < gridSize; i++ {
		var row []image.Image
		for j := 0; j < gridSize; j++ {
			row = append(row, img.(helpers.SubImager).SubImage(
				image.Rect(
					partitionSizeX*j, partitionSizeY*i,
					partitionSizeX*j+partitionSizeX, partitionSizeY*i+partitionSizeY)))
		}
		grid = append(grid, row)
	}

	// make emojis

	var emojiBaseName, exists2 = ctx.Arguments[emojiNameArg]
	if !exists2 {
		var _, err3 = ctx.Session.ChannelMessageSend(ctx.Message.ChannelID, "Emoji name not found!")
		return err3
	}

	var messageString = ""

	for i := 0; i < gridSize; i++ {
		for j := 0; j < gridSize; j++ {

			var buffer bytes.Buffer

			//goland:noinspection GoNilness
			var err4 = png.Encode(&buffer, grid[i][j])
			if err4 != nil {
				return err4
			}

			var encodedImage = base64.StdEncoding.EncodeToString(buffer.Bytes())

			var emoji, err5 = ctx.Session.GuildEmojiCreate(
				ctx.Message.GuildID,
				// humans naturally start counting at 1, not 0
				fmt.Sprint(emojiBaseName, "_", i+1, "_", j+1),
				fmt.Sprint("data:png;base64,", encodedImage),
				nil)
			if err5 != nil {
				return err5
			}

			messageString += fmt.Sprint(":", emoji.Name, ":")
		}
		messageString += "\n"
	}

	var _, err3 = ctx.Session.ChannelMessageSend(
		ctx.Message.ChannelID,
		fmt.Sprint(
			"Copy the text below:\n",
			messageString))

	return err3
}
