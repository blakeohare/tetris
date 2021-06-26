import UIKit

class ViewController: UIViewController {

    /*
    var cluesLabel: UILabel!
    var answersLabel: UILabel!
    var currentAnswer: UITextField!
    var scoreLabel: UILabel!
    var letterButtons = [UIButton]()
    */
    //var x = 0
    // var y = 0
    var game: Game? = nil
    var actions: [String] = []
    
    
    override func viewDidLoad() {
        super.viewDidLoad()
        // Do any additional setup after loading the view.
    }
    
    @objc func pressLeft(_ sender: UIButton) {
        self.actions.append("left")
    }
    
    @objc func pressRight(_ sender: UIButton) {
        self.actions.append("right")
    }
    
    @objc func pressFlip(_ sender: UIButton) {
        self.actions.append("flip")
    }
    
    func doGameLoopIteration(_ surface: DrawSurface) {
        doUpdate(surface)
        DispatchQueue.main.asyncAfter(deadline: .now() + 1.0 / 30) {
            self.doGameLoopIteration(surface)
        }
    }
    
    func doUpdate(_ surface: DrawSurface) {
        if self.game == nil {
            self.game = Game()
        }
        surface.fill(0, 0, 0)
        self.game!.update(actions: self.actions)
        self.actions.removeAll()
        surface.fill(0, 0, 0)
        self.game!.render(g: surface)
        surface.setNeedsDisplay()
    }
    
    override func loadView() {
        view = UIView()
        view.backgroundColor = .red
        
        let leftBtn = UIButton()
        leftBtn.setTitle("<--", for: .normal)
        leftBtn.translatesAutoresizingMaskIntoConstraints = false
        leftBtn.addTarget(self, action: #selector(pressLeft), for: .touchUpInside)
        view.addSubview(leftBtn)
        
        let flipBtn = UIButton()
        flipBtn.setTitle("Flip", for: .normal)
        flipBtn.translatesAutoresizingMaskIntoConstraints = false
        flipBtn.addTarget(self, action: #selector(pressFlip), for: .touchUpInside)
        view.addSubview(flipBtn)
        
        let rightBtn = UIButton()
        rightBtn.setTitle("-->", for: .normal)
        rightBtn.translatesAutoresizingMaskIntoConstraints = false
        rightBtn.addTarget(self, action: #selector(pressRight), for: .touchUpInside)
        view.addSubview(rightBtn)
        
        let ds = DrawSurface()
        ds.translatesAutoresizingMaskIntoConstraints = false
        ds.backgroundColor = .yellow
        view.addSubview(ds)
     
        DispatchQueue.main.asyncAfter(deadline: .now() + 0.1) {
            self.doGameLoopIteration(ds)
        }
        /*
        scoreLabel = UILabel()
        scoreLabel.translatesAutoresizingMaskIntoConstraints = false
        scoreLabel.textAlignment = .right
        scoreLabel.text = "Score: 0"
        view.addSubview(scoreLabel)
        
        cluesLabel = UILabel()
        cluesLabel.translatesAutoresizingMaskIntoConstraints = false
        cluesLabel.font = UIFont.systemFont(ofSize: 24)
        cluesLabel.text = "CLUES"
        cluesLabel.numberOfLines = 0
        cluesLabel.backgroundColor = .red
        view.addSubview(cluesLabel)

        answersLabel = UILabel()
        answersLabel.translatesAutoresizingMaskIntoConstraints = false
        answersLabel.font = UIFont.systemFont(ofSize: 24)
        answersLabel.text = "ANSWERS"
        answersLabel.numberOfLines = 0
        answersLabel.textAlignment = .right
        answersLabel.backgroundColor = .blue
        view.addSubview(answersLabel)
        
        currentAnswer = UITextField()
        currentAnswer.translatesAutoresizingMaskIntoConstraints = false
        currentAnswer.placeholder = "Tap letters to guess"
        currentAnswer.textAlignment = .center
        currentAnswer.font = UIFont.systemFont(ofSize: 44)
        currentAnswer.isUserInteractionEnabled = false
        currentAnswer.backgroundColor = .yellow
        view.addSubview(currentAnswer)
        */
        NSLayoutConstraint.activate([
            leftBtn.bottomAnchor.constraint(equalTo: view.layoutMarginsGuide.bottomAnchor, constant: -30),
            leftBtn.leftAnchor.constraint(equalTo: view.layoutMarginsGuide.leftAnchor),
            rightBtn.bottomAnchor.constraint(equalTo: view.layoutMarginsGuide.bottomAnchor, constant: -30),
            rightBtn.rightAnchor.constraint(equalTo: view.layoutMarginsGuide.rightAnchor),
            flipBtn.bottomAnchor.constraint(equalTo: view.layoutMarginsGuide.bottomAnchor, constant: -30),
            flipBtn.centerXAnchor.constraint(equalTo: view.layoutMarginsGuide.centerXAnchor),
            
            
            ds.leftAnchor.constraint(equalTo: view.layoutMarginsGuide.leftAnchor),
            ds.rightAnchor.constraint(equalTo: view.layoutMarginsGuide.rightAnchor),
            ds.topAnchor.constraint(equalTo: view.layoutMarginsGuide.topAnchor, constant: 100),
            ds.bottomAnchor.constraint(equalTo: view.layoutMarginsGuide.bottomAnchor, constant: -100),
            /*
            scoreLabel.topAnchor.constraint(equalTo: view.layoutMarginsGuide.topAnchor),
            scoreLabel.trailingAnchor.constraint(equalTo: view.layoutMarginsGuide.trailingAnchor),
            
            cluesLabel.topAnchor.constraint(equalTo: scoreLabel.bottomAnchor),
            cluesLabel.leadingAnchor.constraint(equalTo: view.layoutMarginsGuide.leadingAnchor, constant: 10),
            cluesLabel.widthAnchor.constraint(equalTo: view.layoutMarginsGuide.widthAnchor, multiplier: 0.6, constant: -10),

            answersLabel.topAnchor.constraint(equalTo: scoreLabel.bottomAnchor),
            answersLabel.trailingAnchor.constraint(equalTo: view.layoutMarginsGuide.trailingAnchor, constant: -10),
            answersLabel.widthAnchor.constraint(equalTo: view.layoutMarginsGuide.widthAnchor, multiplier: 0.4, constant: -10),
            answersLabel.heightAnchor.constraint(equalTo: cluesLabel.heightAnchor),
            
            currentAnswer.centerXAnchor.constraint(equalTo: view.centerXAnchor),
            currentAnswer.widthAnchor.constraint(equalTo: view.widthAnchor, multiplier: 0.5),
            currentAnswer.topAnchor.constraint(equalTo: cluesLabel.bottomAnchor, constant: 20),
            */
        ])
        
    }

}

