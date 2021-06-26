use sdl2::event::Event;
use sdl2::keyboard::Keycode;
use sdl2::pixels::Color;
use sdl2::rect::Rect;
use sdl2::render::Canvas;
use sdl2::video::Window;
use std::time::{Duration, SystemTime, UNIX_EPOCH};

struct Game {
    row_count: i32,
    col_count: i32,
    board: Vec<Vec<i32>>,
    overlay: Vec<Vec<i32>>,
    overlay_x: i32,
    overlay_y: i32,
    overlay_rotate_is_transpose: bool,
    fall_counter: i32,
    clear_counter: i32,
    clearing_lines: Vec<i32>,
}

impl Game {
    fn new() -> Game {
        let mut game = Game {
            row_count: 20,
            col_count: 10,
            board: Vec::new(),
            overlay: Vec::new(),
            overlay_x: 0,
            overlay_y: 0,
            overlay_rotate_is_transpose: false,
            fall_counter: 0,
            clear_counter: 99999,
            clearing_lines: Vec::new(),
        };

        for _x in 0..game.col_count {
            let mut col = Vec::new();
            for _y in 0..game.row_count {
                col.push(0);
            }
            game.board.push(col);
        }

        game
    }

    fn update(&mut self, actions: &Vec<&str>) {
        if self.is_clear_mode() {
            if self.clear_counter < self.col_count {
                for row in &self.clearing_lines {
                    self.board[self.clear_counter as usize][*row as usize] = 0;
                }
                self.clear_counter += 1;
            } else {
                self.clear_counter = 0;
                if !self.clearing_lines.is_empty() {
                    let mut keep_lines: Vec<i32> = Vec::new();
                    for y in 0..self.row_count {
                        if !self.clearing_lines.contains(&y) {
                            keep_lines.push(y);
                        }
                    }
                    for x in 0..self.col_count {
                        let mut new_col: Vec<i32> = Vec::new();
                        for _ in 0..(self.row_count as usize - &keep_lines.len()) {
                            new_col.push(0);
                        }
                        for y in &keep_lines {
                            new_col.push(self.board[x as usize][*y as usize]);
                        }
                        assert_eq!(new_col.len() as i32, self.row_count);
                        self.board[x as usize] = new_col;
                    }
                    self.clearing_lines.clear();
                }

                self.populate_overlay();
                if !self.is_overlay_valid() {
                    println!("The game must end here.");
                }
            }
        } else {
            for action in actions {
                match *action {
                    "left" | "right" => {
                        let dx = if *action == "left" { -1 } else { 1 };
                        let _ = self.try_move_overlay(dx, 0);
                    },
                    "rot1" | "rot2" => {
                        if self.overlay_rotate_is_transpose {
                            self.transpose_overlay();
                            if !self.is_overlay_valid() {
                                self.transpose_overlay();
                            }
                        } else {
                            let is_cw = *action == "rot1";
                            if is_cw {
                                self.transpose_overlay();
                                self.flip_overlay();
                            } else {
                                self.flip_overlay();
                                self.transpose_overlay();
                            }

                            if !self.is_overlay_valid() {
                                if is_cw {
                                    self.flip_overlay();
                                    self.transpose_overlay();
                                } else {
                                    self.transpose_overlay();
                                    self.flip_overlay();
                                }
                            }
                        }
                    },
                    "drop" => {
                        if !self.try_move_overlay(0, 1) {
                            self.flatten_overlay();
                        }
                    },

                    _ => { }
                }
            }

            self.fall_counter -= 1;
            if self.fall_counter <= 0 {
                self.fall_counter = 30;
                if !self.try_move_overlay(0, 1) {
                    self.flatten_overlay();
                }
            }
        }
    }

    fn populate_overlay(&mut self) {
        let piece: [i32; 16];
        let piece_id: i32 = rand::random::<i32>() % 7 + 1;
        let mut rotate_is_transpose = false;
        match piece_id {
            1 => { 
                piece = [
                    0, 0, 0, 0,
                    0, 1, 1, 0,
                    0, 1, 1, 0,
                    0, 0, 0, 0,
                ];
                rotate_is_transpose = true;
            },
            2 => {
                piece = [
                    0, 1, 0, 0,
                    1, 1, 1, 0,
                    0, 0, 0, 0,
                    0, 0, 0, 0,
                ];
            },
            3 => {
                piece = [
                    0, 1, 0, 0, 
                    0, 1, 0, 0, 
                    0, 1, 0, 0, 
                    0, 1, 0, 0, 
                ];
                rotate_is_transpose = true;
            },
            4 => {
                piece = [
                    0, 2, 0, 0,
                    2, 2, 0, 0,
                    2, 0, 0, 0,
                    0, 0, 0, 0,
                ];
            },
            5 => {
                piece = [
                    0, 2, 0, 0,
                    0, 2, 0, 0,
                    2, 2, 0, 0,
                    0, 0, 0, 0,
                ];
            },
            6 => {
                piece = [
                    0, 3, 0, 0,
                    0, 3, 3, 0,
                    0, 0, 3, 0,
                    0, 0, 0, 0,
                ];
            },
            7 => {
                piece = [
                    0, 3, 0, 0,
                    0, 3, 0, 0,
                    0, 3, 3, 0,
                    0, 0, 0, 0,
                ];
            },
            _ => {
                return
            },
        }

        self.overlay.clear();
        for x in 0..4 {
            let mut col: Vec<i32> = Vec::new();
            for y in 0..4 {
                col.push(piece[x + y * 4]);
            }
            self.overlay.push(col);
        }

        self.overlay_x = (self.col_count - 3) / 2;
        self.overlay_y = 0;
        self.overlay_rotate_is_transpose = rotate_is_transpose;
    }

    fn flatten_overlay(&mut self) {
        assert_eq!(self.is_overlay_valid(), true);

        for y in 0..4 {
            for x in 0..4 {
                let value = self.overlay[x as usize][y as usize];
                if value > 0 {
                    let tx = (x + self.overlay_x) as usize;
                    let ty = (y + self.overlay_y) as usize;
                    self.board[tx][ty] = value;
                }
            }
        }

        self.overlay.clear();

        self.clearing_lines.clear();
        for y in 0..self.row_count {
            let mut all_filled = true;
            for x in 0..self.col_count {
                if self.board[x as usize][y as usize] == 0 {
                    all_filled = false;
                    break
                }
            }
            if all_filled {
                self.clearing_lines.push(y);
            }
        }
    }

    fn transpose_overlay(&mut self) {
        for y in 0..4 {
            for x in y + 1..4 {
                let ux = x as usize;
                let uy = y as usize;
                let t = self.overlay[ux][uy];
                self.overlay[ux][uy] = self.overlay[uy][ux];
                self.overlay[uy][ux] = t;
            }
        }
    }

    fn flip_overlay(&mut self) {
        for y in 0..4 {
            let uy = y as usize;
            let t = self.overlay[0][uy];
            self.overlay[0][uy] = self.overlay[2][uy];
            self.overlay[2][uy] = t;
        }
    }

    fn try_move_overlay(&mut self, dx: i32, dy: i32) -> bool {
        self.overlay_x += dx;
        self.overlay_y += dy;
        if self.is_overlay_valid() {
            return true;
        }
        self.overlay_x -= dx;
        self.overlay_y -= dy;
        return false;
    }

    fn is_clear_mode(&self) -> bool {
        self.overlay.len() == 0
    }

    fn is_overlay_valid(&self) -> bool {
        if self.overlay.len() == 0 {
            return true;
        }

        for y in 0..4 {
            for x in 0..4 {
                let used = self.overlay[x as usize][y as usize] > 0;
                if used {
                    let tx = x + self.overlay_x;
                    let ty = y + self.overlay_y;
                    if tx < 0 || ty < 0 || tx >= self.col_count || ty >= self.row_count {
                        return false;
                    }
                    if self.board[tx as usize][ty as usize] > 0 {
                        return false;
                    }
                }
            }
        }
        return true;
    }

    fn render(&self, canvas: &mut Canvas<Window>) {
        
        let tile_size = 24;
        let screen_width = 800;
        let screen_height = 600;
        let left = (screen_width - tile_size * self.col_count) / 2;
        let top = (screen_height - tile_size * self.row_count) / 2;
        
        draw_rect(canvas, left, top, tile_size * self.col_count, tile_size * self.row_count, 0, 0, 0);

        for y in 0..self.row_count {
            for x in 0..self.col_count {
                let mut value: i32 = self.board[x as usize][y as usize];
                
                // if it's empty, check the overlay for values
                if value == 0 && !self.is_clear_mode() {
                    let ox = x - self.overlay_x;
                    let oy = y - self.overlay_y;
                    if ox >= 0 && ox < 4 && oy >= 0 && oy < 4 {
                        value = self.overlay[ox as usize][oy as usize];
                    }
                }

                if value > 0 {
                    let px = left + x * tile_size;
                    let py = top + y * tile_size;
                    if value == 1 {
                        draw_rect(canvas, px, py, tile_size, tile_size, 255, 255, 255);
                    } else if value == 2 {
                        draw_rect(canvas, px, py, tile_size, tile_size, 0, 128, 255);
                    } else {
                        draw_rect(canvas, px, py, tile_size, tile_size, 0, 200, 50);
                    }
                }
            }
        }



        // draw_rect(canvas, self.temp_x * 20, self.temp_y * 20, 20, 20, 255, 255, 0);
    }
}

fn draw_rect(canvas: &mut Canvas<Window>, x: i32, y: i32, width: i32, height: i32, red: i32, green: i32, blue: i32) {
    canvas.set_draw_color(Color::RGB(red as u8, green as u8, blue as u8));
    let output = canvas.fill_rect(Rect::new(x, y, width as u32, height as u32));
    assert_eq!(output.is_ok(), true);
}

fn unix_time_millis() -> u64 {
    let start = SystemTime::now();
    let since_the_epoch = start
        .duration_since(UNIX_EPOCH)
        .expect("Time went backwards");
    let in_ms = since_the_epoch.as_secs() * 1000 +
            since_the_epoch.subsec_nanos() as u64 / 1_000_000;
    in_ms
}

pub fn main() {
    let sdl_context = sdl2::init().unwrap();
    let video_subsystem = sdl_context.video().unwrap();

    let mut game = Game::new();
    let mut actions: Vec<&str> = Vec::new();

    let window = video_subsystem.window("Tetris but like in Rust", 800, 600)
        .position_centered()
        .build()
        .unwrap();

    let mut canvas = window.into_canvas().build().unwrap();
    let mut down_pressed = false;

    canvas.set_draw_color(Color::RGB(0, 0, 0));
    canvas.clear();
    canvas.present();

    let mut event_pump = sdl_context.event_pump().unwrap();
    let mut frame_time = unix_time_millis();

    'running: loop {
        for event in event_pump.poll_iter() {
            match event {
                Event::Quit {..} => {
                    break 'running
                },
                Event::KeyDown { keycode, ..} => {
                    match keycode.unwrap() {
                        Keycode::Left { .. } => { 
                            actions.push("left");
                        },
                        Keycode::Right { .. } => {
                            actions.push("right");
                        },
                        Keycode::Up { .. } => {
                            actions.push("rot1");
                        },
                        Keycode::Down { .. } => {
                            down_pressed = true;
                        },
                        Keycode::Space { .. } => {
                            actions.push("rot2");
                        },
                        Keycode::Escape { .. } => {
                            break 'running
                        },
                        _ => {}
                    }
                },
                Event::KeyUp { keycode, ..} => {
                    if keycode.unwrap() == Keycode::Down {
                        down_pressed = false;
                    }
                },
                _ => {}
            }
        }

        if down_pressed {
            actions.push("drop");
        }
        
        game.update(&actions);
        actions.clear();
        // i = (i + 1) % 255;

        canvas.set_draw_color(Color::RGB(80, 80, 80));
        canvas.clear();
        game.render(&mut canvas);

        canvas.present();
        let now = unix_time_millis();
        let diff = now - frame_time;
        if diff < 16 {
            let wait_time: u32 = (16 - diff) as u32; //.try_into().unwrap();
            let delay32 = wait_time * 1_000_000;
            ::std::thread::sleep(Duration::new(0, delay32));
        }
        frame_time = unix_time_millis();
    }
}