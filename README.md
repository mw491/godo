# __GO__(lang)__DO__(it)
##### A simple todo app written in go, (go)ncurses and objectbox.
![Preview of the app (godo)](screenshots/preview.png "godo preview")



to install objectboxlib to the program:   

```
mkdir objectboxlib && cd objectboxlib
bash <(curl -s https://raw.githubusercontent.com/objectbox/objectbox-c/main/download.sh) 0.13.0
```    

for a more detailed explanation, visit https://golang.objectbox.io/install.   



| Key     | Action               |
| ------- | -------------------- |
| `space` | Toggle Check           | 
| `j`     | Move cursor up         |
| `k`     | Move cursor down       |
| `o`     | Add New Todo           |
| `d`     | Delete Selected Todo   |
| `q`     | Quit App               |
