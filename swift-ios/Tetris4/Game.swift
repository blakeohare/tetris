import Foundation

class Game {
    
    var board: [[Int]]
    var overlay: [[Int]]? = nil
    var gameCounter = 0
    var fallCounter = 0
    var overlayX = 0
    var overlayY = 0
    var overlayRotateIsTranspose = false
    let rowCount = 20
    let colCount = 10
    var mode = "FALL"
    var clearingLines: [Int] = []
    var clearCounter = 0
    var linesCleared = 0
    var level = 0
    
    init() {
        self.board = makeGrid(self.colCount, self.rowCount)
    }
    
    func update(actions: [String]) {
        
        self.gameCounter += 1
        let fallMax = max(30 - self.level, 5)
        
        for action in actions {
            switch action {
                case "left":
                    self.tryMove(-1, 0)
                case "right":
                    self.tryMove(1, 0)
                case "flip":
                    self.rotateOverlay(true)
                default:
                    break
            }
        }
        switch self.mode {
        
            case "FALL":
                if self.overlay == nil {
                    let pieceType = Int.random(in: 1...7)
                    self.overlay = self.createNewOverlay(pieceType: pieceType)
                    self.overlayRotateIsTranspose = pieceType == 1 || pieceType == 2 || pieceType == 5 || pieceType == 7
                    self.fallCounter = 15
                    self.overlayX = 4
                    self.overlayY = 0
                }
                
                self.fallCounter -= 1
                if self.fallCounter <= 0 {
                    self.fallCounter = 15
                    if !self.tryMove(0, 1) {
                        self.flattenOverlay()
                        self.overlay = nil
                        self.checkForClear()
                    }
                }
                
            case "CLEAR":
                for line in self.clearingLines {
                    if self.clearCounter < self.colCount {
                        self.board[self.clearCounter][line] = 0
                    } else {
                        self.clearCounter = 0
                        self.mode = "FALL"
                        self.overlay = nil
                        
                        var newCols: [[Int]] = []
                        for x in 0..<self.colCount {
                            var col: [Int] = []
                            for _ in 0..<self.clearingLines.count {
                                col.append(0)
                            }
                            for y in 0..<self.rowCount {
                                if !self.clearingLines.contains(y) {
                                    col.append(self.board[x][y])
                                }
                            }
                            newCols.append(col)
                        }
                        self.board = newCols
                    }
                }
                self.clearCounter += 1
                
            default:
                break
        }
    }
    
    func checkForClear() {
        var clearRows: [Int] = []
        for y in 0..<self.rowCount {
            var allFilled = true
            for x in 0..<self.colCount {
                if self.board[x][y] == 0 {
                    allFilled = false
                    break
                }
            }
            if allFilled {
                clearRows.append(y)
            }
        }
        self.clearingLines = clearRows
        if self.clearingLines.count > 0 {
            self.linesCleared += clearRows.count
            self.level = self.linesCleared / 10
            self.mode = "CLEAR"
            self.clearCounter = 0
        }
    }
    
    func flattenOverlay() {
        for y in 0..<4 {
            for x in 0..<4 {
                let value = self.overlay![x][y]
                if value > 0 {
                    let tx = x + self.overlayX
                    let ty = y + self.overlayY
                    self.board[tx][ty] = value
                }
            }
        }
    }
    
    func createNewOverlay(pieceType: Int) -> [[Int]] {
        switch pieceType {
            case 1:
                return [
                    [0, 0, 0, 0],
                    [0, 1, 1, 0],
                    [0, 1, 1, 0],
                    [0, 0, 0, 0],
               ]
            case 2:
                return [
                    [0, 1, 0, 0],
                    [0, 1, 0, 0],
                    [0, 1, 0, 0],
                    [0, 1, 0, 0],
               ]
            case 3:
                return [
                    [0, 1, 0, 0],
                    [1, 1, 1, 0],
                    [0, 0, 0, 0],
                    [0, 0, 0, 0],
               ]
            case 4:
                return [
                    [0, 2, 0, 0],
                    [0, 2, 0, 0],
                    [2, 2, 0, 0],
                    [0, 0, 0, 0],
               ]
            case 5:
                return [
                    [0, 2, 0, 0],
                    [2, 2, 0, 0],
                    [2, 0, 0, 0],
                    [0, 0, 0, 0],
               ]
            case 6:
                return [
                    [0, 3, 0, 0],
                    [0, 3, 0, 0],
                    [0, 3, 3, 0],
                    [0, 0, 0, 0],
               ]
            case 7:
                return [
                    [0, 3, 0, 0],
                    [0, 3, 3, 0],
                    [0, 0, 3, 0],
                    [0, 0, 0, 0],
               ]
            default: return [[]]
        }
    }
    
    @discardableResult func rotateOverlay(_ clockwise: Bool) -> Bool {
        var t1 = false
        var flip = false
        var t2 = false
        if self.overlayRotateIsTranspose {
            t1 = true
        } else if clockwise {
            t1 = true
            flip = true
        } else {
            flip = true
            t2 = true
        }
        
        if t1 { self.transposeOverlay() }
        if flip { self.flipOverlay() }
        if t2 { self.transposeOverlay() }
        
        if !self.isOverlayPossible() {
            if t2 { self.transposeOverlay() }
            if flip { self.flipOverlay() }
            if t1 { self.transposeOverlay() }
            return false
        }
        return true
    }
    
    func transposeOverlay() {
        if (self.overlay == nil) { return }
        for y in 0..<4 {
            for x in (y+1)..<4 {
                let t = self.overlay![y][x]
                self.overlay![y][x] = self.overlay![x][y]
                self.overlay![x][y] = t
            }
        }
    }
    
    func flipOverlay() {
        for y in 0..<4 {
            let t = self.overlay![0][y]
            self.overlay![0][y] = self.overlay![2][y]
            self.overlay![2][y] = t
        }
    }
    
    @discardableResult func tryMove(_ dx: Int, _ dy: Int) -> Bool {
        if self.overlay == nil { return false }
        self.overlayX += dx
        self.overlayY += dy
        if !self.isOverlayPossible() {
            self.overlayX -= dx
            self.overlayY -= dy
            return false
        }
        return true
    }

    func isOverlayPossible() -> Bool {
        for oy in 0..<4 {
            for ox in 0..<4 {
                let value = self.overlay![ox][oy]
                if value > 0 {
                    let x = ox + self.overlayX
                    let y = oy + self.overlayY
                    if x < 0 || y < 0 { return false }
                    if x >= self.colCount || y >= self.rowCount { return false }
                    if self.board[x][y] > 0 { return false }
                }
            }
        }
        return true
    }
    
    func render(g: DrawSurface) {
        // g.drawRect(self.x, self.y, 10, 10, 255, 255, 0)
        let tileSize = 14
        let left = 70
        let top = 30
        
        let themes : [[[Int]]] = [
            [[0, 128, 255], [0, 200, 80]],
            [[255, 255, 0], [255, 128, 0]],
        ]
        
        let themeSlice = themes[self.level % themes.count]
        let theme : [[Int]?] = [nil, [255, 255, 255], themeSlice[0], themeSlice[1]]
        
        g.drawRect(left, top, tileSize * self.colCount, tileSize * self.rowCount, 50, 50, 50)
        
        for x in 0..<self.colCount {
            for y in 0..<self.rowCount {
                var value = 0
                if self.overlay != nil {
                    let ox = x - self.overlayX
                    let oy = y - self.overlayY
                    if ox >= 0 && oy >= 0 && ox < 4 && oy < 4 {
                        value = self.overlay![ox][oy]
                    }
                }
                if value == 0 {
                    value = self.board[x][y]
                }
                if value > 0 {
                    let color = theme[value]
                    if color != nil {
                        g.drawRect(left + x * tileSize, top + y * tileSize, tileSize, tileSize, color![0], color![1], color![2])
                    }
                }
            }
        }
        
    }
}
