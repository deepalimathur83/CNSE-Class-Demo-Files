/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"errors"
	"fmt"
	"os"

	"drexel.edu/todo/db"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

var (
	dbFileNameFlag string
	restoreDbFlag  bool
	listFlag       bool
	itemStatusFlag bool
	queryFlag      int
	addFlag        string
	updateFlag     string
	deleteFlag     int
)

type AppOptType int

const (
	LIST_DB_ITEM AppOptType = iota
	RESTORE_DB_ITEM
	QUERY_DB_ITEM
	ADD_DB_ITEM
	UPDATE_DB_ITEM
	DELETE_DB_ITEM
	CHANGE_ITEM_STATUS
	NOT_IMPLEMENTED
	INVALID_APP_OPT
)

func processCmdLineFlags(cmd *cobra.Command) (AppOptType, error) {

	var appOpt AppOptType = INVALID_APP_OPT

	//show help if no flags are set
	if len(os.Args) == 1 {
		cmd.Flags().Usage()
		return appOpt, errors.New("no flags were set")
	}

	// Loop over the flags and check which ones are set, set appOpt
	// accordingly
	cmd.Flags().Visit(func(f *pflag.Flag) {
		switch f.Name {
		case "list":
			appOpt = LIST_DB_ITEM
		case "restore":
			appOpt = RESTORE_DB_ITEM
		case "query":
			appOpt = QUERY_DB_ITEM
		case "add":
			appOpt = ADD_DB_ITEM
		case "upadate":
			appOpt = UPDATE_DB_ITEM
		case "delete":
			appOpt = DELETE_DB_ITEM

		//TODO: EXTRA CREDIT - Implment the -s flag that changes the
		//done status of an item in the database.  For example -s=true
		//will set the done status for a particular item to true, and
		//-s=false will set the done states for a particular item to
		//false.
		//
		//HINT FOR EXTRA CREDIT
		//Note the -s option also requires an id for the item to that
		//you want to change.  I recommend you use the -q option to
		//specify the item id.  Therefore, the -s option is only valid
		//if the -q option is also set
		case "s":
			//For extra credit you will need to change some things here
			//and also in main under the CHANGE_ITEM_STATUS case

			if appOpt == QUERY_DB_ITEM {
				appOpt = CHANGE_ITEM_STATUS
			} else {
				appOpt = INVALID_APP_OPT
			}

		default:
			appOpt = INVALID_APP_OPT
		}
	})

	if appOpt == INVALID_APP_OPT || appOpt == NOT_IMPLEMENTED {
		fmt.Println("Invalid option set or the desired option is not currently implemented")
		cmd.Help()
		return appOpt, errors.New("no flags or unimplemented were set")
	}

	return appOpt, nil
}

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "todo",
	Short: "A brief description of your application",
	Long: `A longer description that spans multiple lines and likely contains
examples and usage of using your application. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	Run: func(cmd *cobra.Command, args []string) {

		opts, err := processCmdLineFlags(cmd)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		//Create a new db object
		todo, err := db.New(dbFileNameFlag)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		//Switch over the command line flags and call the appropriate
		//function in the db package
		switch opts {
		case RESTORE_DB_ITEM:
			fmt.Println("Running RESTORE_DB_ITEM...")
			if err := todo.RestoreDB(); err != nil {
				fmt.Println("Error: ", err)
				break
			}
			fmt.Println("Database restored from backup file")
		case LIST_DB_ITEM:
			fmt.Println("Running QUERY_DB_ITEM...")
			todoList, err := todo.GetAllItems()
			if err != nil {
				fmt.Println("Error: ", err)
				break
			}
			for _, item := range todoList {
				todo.PrintItem(item)
			}
			fmt.Println("THERE ARE", len(todoList), "ITEMS IN THE DB")
			fmt.Println("Ok")

		case QUERY_DB_ITEM:
			fmt.Println("Running QUERY_DB_ITEM...")
			item, err := todo.GetItem(queryFlag)
			if err != nil {
				fmt.Println("Error: ", err)
				break
			}
			todo.PrintItem(item)
			fmt.Println("Ok")
		case ADD_DB_ITEM:
			fmt.Println("Running ADD_DB_ITEM...")
			item, err := todo.JsonToItem(addFlag)
			if err != nil {
				fmt.Println("Add option requires a valid JSON todo item string")
				fmt.Println("Error: ", err)
				break
			}
			if err := todo.AddItem(item); err != nil {
				fmt.Println("Error: ", err)
				break
			}
			fmt.Println("Ok")
		case UPDATE_DB_ITEM:
			fmt.Println("Running UPDATE_DB_ITEM...")
			item, err := todo.JsonToItem(updateFlag)
			if err != nil {
				fmt.Println("Update option requires a valid JSON todo item string")
				fmt.Println("Error: ", err)
				break
			}
			if err := todo.UpdateItem(item); err != nil {
				fmt.Println("Error: ", err)
				break
			}
			fmt.Println("Ok")
		case DELETE_DB_ITEM:
			fmt.Println("Running DELETE_DB_ITEM...")
			err := todo.DeleteItem(deleteFlag)
			if err != nil {
				fmt.Println("Error: ", err)
				break
			}
			fmt.Println("Ok")

		case CHANGE_ITEM_STATUS:

			//For the CHANGE_ITEM_STATUS extra credit you will also
			//need to add some code here

			fmt.Println("Running CHANGE_ITEM_STATUS...")

			err := todo.ChangeItemDoneStatus(queryFlag, itemStatusFlag)

			if err != nil {

				fmt.Println("Error: ", err)

				break
			}
			fmt.Println("OK")

			//fmt.Println("Not implemented yet, but it can be for extra credit")

		default:
			fmt.Println("INVALID_APP_OPT")
		}
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	// rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.todo.yaml)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.

	rootCmd.PersistentFlags().StringVarP(&addFlag, "add", "a", "", "Adds an item to the database")
	rootCmd.PersistentFlags().StringVar(&dbFileNameFlag, "db", "./data/todo.json", "Name of the target database file (default: \"./data/todo.json\")")
	rootCmd.PersistentFlags().IntVarP(&deleteFlag, "delete", "d", 0, "Deletes an item from the database")
	rootCmd.PersistentFlags().BoolVarP(&listFlag, "list", "l", false, "List all the items in the database")
	rootCmd.PersistentFlags().IntVarP(&queryFlag, "query", "q", 0, "Query an item in the database")
	rootCmd.PersistentFlags().BoolVarP(&restoreDbFlag, "restore", "r", false, "Restore the database from the backup file")
	rootCmd.PersistentFlags().BoolVarP(&itemStatusFlag, "status", "s", false, "Change item 'done' status to true or false. must be used alongside -q flag. ex: -q<target:int> -s=<isDone:bool>")
	rootCmd.PersistentFlags().BoolP("toggle", "t", false, "Help message for toggle")
	rootCmd.PersistentFlags().StringVarP(&updateFlag, "update", "u", "", "Updates an item in the database")

}
