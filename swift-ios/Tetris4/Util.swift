import Foundation

func makeGrid(_ width: Int, _ height: Int) -> [[Int]] {
    var cols: [[Int]] = []
    for _ in 0..<width {
        var col: [Int] = []
        for _ in 0..<height {
            col.append(0)
        }
        cols.append(col)
    }
    return cols
}
