package commands

import "fmt"

const usageInfo = `  __________________________________
    _______             _
   |__   __|  ___      | |   ___
      | |    /   \   __| |  /   \
      | |    |[ ]|  | _  |  |[ ]|
      |_|    \___/  \____|  \___/
  __________________________________
`

func Default() {
	fmt.Print(usageInfo)
}