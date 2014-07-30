import "fmt"
import "os"

func isValid(board, x, y, c) {
	for i = 0; i < 9; i++ {
		if board[x][i] == c {
			return false
		}
		if board[i][y] == c {
			return false
		}
	}
	for i = 3*(x/3); i < 3*(x/3+1); i++ {
		for j = 3*(y/3); j < 3*(y/3+1); j++ {
			if board[i][j] == c {
				return false
			}
		}
	}
	return true
}

count = 0

func showBoard(board) {
	fmt.Println("-----------------------------------")
	for i, line = range board {
		fmt.Println(line)
	}
}

func solveSudoku(board) {
	showBoard(board)
	count++
	fmt.Println(count)

	for i = 0; i < board.Length(); i++ {
		for j = 0; j < board[i].Length(); j++ {
			if board[i][j] == "." {
				for k = 0; k < 9; k++ {
					c = "" + (k+1)
					if isValid(board, i, j, c) {
						board[i][j] = c
						if solveSudoku(board) {
							return true
						}
						board[i][j] = "."
					}
				}
				return false
			}
		}
	}
	return true
}

theboard = [
	['5','3','.','.','7','.','9','.','.'],
	['6','.','.','1','9','5','.','.','.'],
	['.','9','8','.','.','.','.','6','.'],
	['8','.','.','.','6','.','.','.','3'],
	['4','.','6','8','.','3','7','.','1'],
	['7','.','.','.','2','.','.','.','6'],
	['.','6','1','.','.','.','2','8','.'],
	['.','.','.','4','1','9','.','.','5'],
	['3','.','5','.','8','.','.','7','9']
]

solveSudoku(theboard)
