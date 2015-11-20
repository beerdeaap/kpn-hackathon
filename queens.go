package main

import (
    "time"
    "math/rand"
    "encoding/json"
    "fmt"
    "net/http"
)

type Position struct {
    X int
    Y int
}

type Board struct {
    Columns, Rows, QueenCount, MaxQueensOnSight int
    InitialQueens, Queens []Position
    Grid [][]bool
    QueensOnSight [][]int
    CandidatePositions []Position
    Fitness int
    AddedQueens []Position
}

type IncomingMessage struct {
    Columns, Rows, Max_queens_on_sight float64
    Initial_queens []map[string]float64
}

type OutgoingMessage struct {
    Added_queens []Position
}

func (board *Board) InitializeGrid() {
    board.Grid = make([][]bool, board.Rows)
    for row := 0; row < board.Rows; row++ {
        board.Grid[row] = make([]bool, board.Columns)
        for col := 0; col < board.Columns; col++{
            board.Grid[row][col] = false
        }
    }
}

func (board *Board) InitializeQueensOnSight() {
    board.QueensOnSight = make([][]int, board.Rows)
    for row := 0; row < board.Rows; row++ {
        board.QueensOnSight[row] = make([]int, board.Columns)
        for col := 0; col < board.Columns; col++{
            board.QueensOnSight[row][col] = 0
        }
    }
}


func (board *Board) QueensToGrid(){
    for _, queen := range board.Queens {
        board.Grid[queen.X][queen.Y] = true
    }
}

func (board *Board) HasQueen(x int, y int) bool {
    return board.Grid[x][y]
}


func (board *Board) AddQueens(queens []Position){
    board.Queens = make([]Position, len(queens))
    copy(board.Queens, queens)
    board.InitializeGrid()
    board.InitializeQueensOnSight()
    board.QueensToGrid()
    board.UpdateQueensOnSight()
    board.Fitness = board.GetFitness()
}

func (board *Board) Initalize(message IncomingMessage) {

    //fmt.Println(message)
    board.Rows = int(message.Rows)
    board.Columns = int(message.Columns)
    board.MaxQueensOnSight = int(message.Max_queens_on_sight)

    for _, position := range message.Initial_queens {
        var pos Position
        pos.X = int(position["x"])
        pos.Y = int(position["y"])
        board.InitialQueens = append(board.InitialQueens, pos)
    }

    board.AddQueens(board.InitialQueens)
    board.AddedQueens = make([]Position, 0)
}

func (board *Board) FindQueens (start_x int, start_y int, dir_x int, dir_y int) int {
    // finds amount of queens in a certain direction
    total := 0
    row := start_x
    col:= start_y
    i := 0
    for  ;row >= 0 && row < board.Rows && col >= 0 && col < board.Columns ;{
        if i != 0 && board.Grid[row][col] == true {
            total++
        }
        row+=dir_x
        col+=dir_y
        i++
    }

    return total
}

func (board *Board) UpdateQueensOnSight (){
    //for every position  (queen or no queen) find out the queens on sight
     // FIXME: can probably done in a go routine
     board.CandidatePositions = make([]Position, 0)
    for row := 0; row < board.Rows; row++ {
        for col := 0; col < board.Columns; col++ {
            queensOnSight := 0
            //North
            queensOnSight += board.FindQueens(row, col, -1, 0)
            //NorthEast
            queensOnSight += board.FindQueens(row, col, -1, 1)
            //East
            queensOnSight += board.FindQueens(row, col, 0, 1)
            //SouthEast
            queensOnSight += board.FindQueens(row, col, 1, 1)
            //South
            queensOnSight += board.FindQueens(row, col, 1, 0)
            //SouthWest
            queensOnSight += board.FindQueens(row, col, 1, -1)
            //West
            queensOnSight += board.FindQueens(row, col, 0, -1)
            //NorthWest
            queensOnSight += board.FindQueens(row, col, -1, -1)
            board.QueensOnSight[row][col] = queensOnSight

            if queensOnSight <= board.MaxQueensOnSight  && board.Grid[row][col] == false{
                candidatePosition := Position{X:row, Y:col}
                board.CandidatePositions = append(board.CandidatePositions, candidatePosition)
                // FIXME: sort these
            }
        }
    }
}

func (board *Board) RemoveQueen(position Position) {

    // Find a queen to remove
    // FIXME: Better to remove the queen with most others on sight
    rand.Seed(time.Now().Unix())
    //toRemove := board.AddedQueens[rand.Intn(len(board.AddedQueens) - 0) + 0]

    // remove from queens
    // remove from added queens
}


func (board *Board) AddQueen(position Position) {
    // Add queen to board, updates grid, queens on sight
    board.Queens = append(board.Queens, position)
    board.AddedQueens = append(board.AddedQueens, position)
    board.QueensToGrid()
    board.UpdateQueensOnSight()
    board.Fitness = board.GetFitness()
}

func (board *Board) GetFitness() int {
    queen_count := 0
    for row := 0; row < board.Rows; row++ {
        for col := 0; col < board.Columns; col++ {
            if board.Grid[row][col]  == true {
                if board.QueensOnSight[row][col] > board.MaxQueensOnSight {
                    return 0
                } else {
                    queen_count++
                }
            }
        }
    }

    return queen_count
}

func (board *Board) Display() {
    for row := 0; row < board.Rows; row++ {
        for col := 0; col < board.Columns; col++ {
            if board.Grid[row][col] == true {
                fmt.Printf("Q ")
            } else {
                fmt.Printf(". ")
            }
        }
        fmt.Println("")
    }
}

func (board *Board) GetCandidatePosition() Position {
    // Return a new random position where qwe can place a queen
    rnd := rand.Intn(len(board.CandidatePositions))
    //fmt.Printf("idx: %d", rnd)
    //fmt.Printf("len: %d\n", len(board.CandidatePositions))
    //fmt.Printf("pos: x, y")
    return board.CandidatePositions[rnd]
}

func (board *Board) DisplayQueensOnSight() {
    for row := 0; row < board.Rows; row++ {
        for col := 0; col < board.Columns; col++ {
            fmt.Printf("%d\t", board.QueensOnSight[row][col])
        }
        fmt.Println("")
    }
}

func (board *Board) FindBest(best *[]Position, message IncomingMessage) {
    // FIXME: add a timer
    bestFitness := 0
    fitnessRepeat := 0
    lastFitness := 0
    max := 1000

    if board.Rows * board.Columns > 2000 {
        max = 100
    } else {
        max = 1000
    }

    for i:=0; i < max; i++ {
        if len(board.CandidatePositions) > 0 {
            board.AddQueen(board.GetCandidatePosition())
            if board.Fitness > bestFitness {
                // board.Display()
                //fmt.Println("")
                *best = make([]Position, len(board.AddedQueens))
                copy(*best, board.AddedQueens)
                bestFitness = board.Fitness
            }
        } else if fitnessRepeat > 10 {
            board.Initalize(message)
        } else {
            board.Initalize(message)
        }

        if board.Fitness == lastFitness{
            fitnessRepeat++
        }
        lastFitness=board.Fitness
    }
}


func MaxQueensHandler(writer http.ResponseWriter, request *http.Request) {

    rand.Seed(time.Now().Unix())
    // handle incoming json
    var message IncomingMessage
    json_blob := make([]byte, request.ContentLength)
    request.Body.Read(json_blob)
    json.Unmarshal(json_blob, &message)
    // fmt.Println(string(json_blob))
    // fmt.Println(message)

    var board Board
    board.Initalize(message)

    bestPositions := make([]Position, 0)
    board.FindBest(&bestPositions, message)

    board.AddQueens(bestPositions)
    // board.Display()
    // board.DisplayQueensOnSight()

    // Send back json
    var out OutgoingMessage
    out.Added_queens = bestPositions
    bytes, _ := json.Marshal(out)
    writer.Write(bytes)
}


func main () {
    http.HandleFunc("/max_queens", MaxQueensHandler)
    http.ListenAndServe(":8080", nil)
}
