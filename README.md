# GopherStrike
GopherStrike is funny 2D multiplayer shooter with physics. Have fun!

<img width="600" alt="image" src="https://github.com/user-attachments/assets/196cd7a9-8819-4e0f-a233-32d2e97ef1e3">



## How to start
Currently, only the local mode is available:
1. Someone has to start the server:
```
go run github.com/simpletonDL/GoGames/server/main@v0.1.6 <port>
```
e.g.
```
go run github.com/simpletonDL/GoGames/server/main@v0.1.6 5006
```
2. Then everyone can connect over a local network using client:
```
go run github.com/simpletonDL/GoGames/client/main@v0.1.6 <server ip> <server port> <nickname>
```
e.g.
```
go run github.com/simpletonDL/GoGames/client/main@v0.1.6 192.168.226.164 5006 awesome.guy
```

## How to play

### Select team
After connecting to the server you appear on the place for command selection. To choose your command (red or blue) just take either the left or right side of the platform.
When everyone has chosen a team press enter to confirm your choise (in such case you become immovable). The game will start when everyone confirms the choice. 
If you want to cancel your choise and become movable press escape.

<img width="600" alt="image" src="https://github.com/user-attachments/assets/6b4c1a42-a669-4892-9eb7-2839dc1e5d58">

### Win the game
After everyone has chosen their team, the GopherStrike starts! Unlike other shooters, the weapon **does not cause any damage**. It just pushes players away. 
So the main goal is to **push enemies off the platform**.
Keyboard controls are intuitive:
* A,← and D,→: move left/right
* W,↑: Jump/Double jump
* S,↓: Move down (in the air) / drop down through platform
* Space: shoot

Also during the game session boxes of weapons spawn:

<img width="150" alt="image" src="https://github.com/user-attachments/assets/e4922c93-a286-4bb1-9c68-a62ac6bc9817">

To get a weapon you just need to touch these boxes. Each weapon has its own total bullets cout, magazie size, recoil (be careful with sniper rifle) and recharge time.

In the upper left/right corner there is a info about your (and others) lives and avalible bullets:

<img width="465" alt="image" src="https://github.com/user-attachments/assets/39804fef-2831-4f1a-b09a-aec0cbc45405">

Format is the following: `nickname L=lives B=bullets in magazine/magazine capacity (total remaining bullets count)`

If all players of some team are dead, then another team wins and everybody are returned back to select team stage.
