import UIKit
import CoreGraphics

class DrawSurface: UIView {
    
    var commandTypes: [String] = ["F"]
    var commandArgs: [Int] = [0, 0, 0]
    
    override func draw(_ rect: CGRect) {
        let _ctx = UIGraphicsGetCurrentContext()
        if (_ctx == nil) {
            print("No graphics context")
            return
        }
        let ctx = _ctx!
        
        var commandIndex = 0
        var argIndex = 0
        var r = 0
        var g = 0
        var b = 0
        var a = 0
        var x = 0
        var y = 0
        var w = 0
        var h = 0
        
        while commandIndex < commandTypes.count {
            let cmd = commandTypes[commandIndex]
            commandIndex += 1
            switch cmd {
                case "F":
                    r = commandArgs[argIndex]
                    g = commandArgs[argIndex + 1]
                    b = commandArgs[argIndex + 2]
                    argIndex += 3
                    self._drawRectImpl(0, 0, Int(rect.width), Int(rect.height), r, g, b, 255)
                case "R":
                    x = commandArgs[argIndex]
                    y = commandArgs[argIndex + 1]
                    w = commandArgs[argIndex + 2]
                    h = commandArgs[argIndex + 3]
                    r = commandArgs[argIndex + 4]
                    g = commandArgs[argIndex + 5]
                    b = commandArgs[argIndex + 6]
                    a = commandArgs[argIndex + 7]
                    self._drawRectImpl(x, y, w, h, r, g, b, a)
                    argIndex += 8
                default:
                    print("Unknown command: \(cmd)")
            }
        }
        
        ctx.flush()
    }
    
    func fill(_ red: Int, _ green: Int, _ blue: Int) {
        commandTypes.removeAll()
        commandArgs.removeAll()
        
        commandTypes.append("F")
        commandArgs.append(red)
        commandArgs.append(green)
        commandArgs.append(blue)
    }
    
    func drawRectAlpha(_ left: Int, _ top: Int, _ width: Int, _ height: Int, _ red: Int, _ green: Int, _ blue: Int, _ alpha: Int) {
        commandTypes.append("R")
        commandArgs.append(left)
        commandArgs.append(top)
        commandArgs.append(width)
        commandArgs.append(height)
        commandArgs.append(red)
        commandArgs.append(green)
        commandArgs.append(blue)
        commandArgs.append(alpha)
    }
    
    func drawRect(_ left: Int, _ top: Int, _ width: Int, _ height: Int, _ red: Int, _ green: Int, _ blue: Int) {
        drawRectAlpha(left, top, width, height, red, green, blue, 255)
    }
    
    func _drawRectImpl(_ left: Int, _ top: Int, _ width: Int, _ height: Int, _ red: Int, _ green: Int, _ blue: Int, _ alpha: Int) {
        let ctx = UIGraphicsGetCurrentContext()
        ctx!.setFillColor(
            red: CGFloat(Float(red) / 255.0),
            green: CGFloat(Float(green) / 255.0),
            blue: CGFloat(Float(blue) / 255),
            alpha: CGFloat(Float(alpha) / 255.0))
        ctx!.addRect(CGRect(x: left, y: top, width: width, height: height))
        ctx!.drawPath(using: .fill)
    }
}
