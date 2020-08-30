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
	"strconv"

	_ "golang.org/x/image/webp"

	"rpdSix/commands"
	"rpdSix/helpers"
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
	emojiNameArg    = "emojiName"
	gridSizeArg     = "gridSize"
	defaultGridSize = 2
)

func run(ctx commands.CommandContext) error {
	// Check if the author has the manage emojis permission
	var messageGuild, err1 = ctx.Session.Guild(ctx.Message.GuildID)
	if err1 != nil {
		return err1
	}
	for _, member := range messageGuild.Members {
		if member.User.ID == ctx.Message.Author.ID {
			for _, roleID := range member.Roles {
				for _, guildRole := range messageGuild.Roles {
					if guildRole.ID == roleID {
						// MANAGE_EMOJIS
						if guildRole.Permissions&0x40000000 == 0x40000000 {
							goto successfulCheck
						}
					}
				}
			}
			goto failedCheck
		}
	}
	return errors.New(fmt.Sprint("member not found, member id: ", ctx.Message.Author.ID))

failedCheck:
	if _, err2 := ctx.Message.Reply("You do not have the manage emojis permission!");
		true {
		return err2
	}

successfulCheck:
	if len(ctx.Message.Attachments) == 0 {
		var _, err3 = ctx.Message.Reply("No attachments were found!")
		return err3
	}

	var attachment = ctx.Message.Attachments[0]
	// consider using ProxyURL instead of URL
	resp, err4 := ctx.Session.Client.Get(attachment.URL)
	if err4 != nil {
		return err4
	} else if resp.StatusCode != 200 {
		return errors.New(fmt.Sprint("status code is not 200, status code is: ", resp.StatusCode))
	}

	var img, _, err5 = image.Decode(resp.Body)
	if err5 != nil {
		return err5
	}

	var gridSize int

	var gridSizeStr, exists2 = ctx.Arguments[gridSizeArg]

	if !exists2 {
		gridSize = defaultGridSize
	} else {
		var gridSizeTmp, err6 = strconv.Atoi(gridSizeStr)
		if err6 != nil {
			return err6
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

	var emojiBaseName, exists3 = ctx.Arguments[emojiNameArg]
	if !exists3 {
		var _, err7 = ctx.Message.Reply("Emoji name not found!")
		return err7
	}

	var messageString = ""

	for i := 0; i < gridSize; i++ {
		for j := 0; j < gridSize; j++ {

			var buffer bytes.Buffer

			//goland:noinspection GoNilness
			var err8 = png.Encode(&buffer, grid[i][j])
			if err8 != nil {
				return err8
			}

			var encodedImage = base64.StdEncoding.EncodeToString(buffer.Bytes())

			var emoji, err9 = ctx.Session.GuildEmojiCreate(
				ctx.Message.GuildID,
				// humans naturally start counting at 1, not 0
				fmt.Sprint(emojiBaseName, "_", i+1, "_", j+1),
				fmt.Sprint("data:png;base64,", encodedImage),
				nil)
			if err9 != nil {
				return err9
			}

			messageString += fmt.Sprint(":", emoji.Name, ":")
		}
		messageString += "\n"
	}

	var _, err10 = ctx.Message.Reply(
		fmt.Sprint(
			"Copy the text below:\n",
			messageString))

	return err10
}
