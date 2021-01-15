# Homeworks

>*– Che cosa sia questa musica?*   
>*Peccato che io debba andare a scuola, se no...*   
>*E rimase lì perplesso. A ogni modo, bisognava prendere una risoluzione: o a scuola, o a sentire i pifferi.*   
>*– Oggi anderò a sentire i pifferi, e domani a scuola: per andare a scuola c’è sempre tempo, – disse finalmente quel monello facendo una spallucciata.*   

>[Carlo Collodi, Le avventure di Pinocchio](http://www.letteraturaitaliana.net/pdf/Volume_9/t217.pdf)


On my way to becoming what I was meant to be, I got diverted by little things so many times that it became an habit.    
On my way to mastering Python, I got enchanted by [micro](https://micropython.org/) and [circuit](https://circuitpython.org/), the two buddies from the hot board.   
On my way to Go, even then, I stumbled upon some [tiny](https://tinygo.org/) distractions.   
But in any case, I must tell, I had way more fun than expected.   

I'm sharing these **homeworks** because there's always time to get back to school once the music is over.     
Unless, by the way, You're looking for a brand new school book...  


## Changelog

* Initial import: examples of Tinygo running on an Arduino Nano33 Iot

I'm looking for a way to control an Arduino Nano33 Iot over Mqtt.   
I'd like to use the board to get some data from a sensor and drive a Neopixel strip.   
I did translate the Neopixel Rainbow animation from a Python example into Go/Tinygo in the past. 
It used to run on an Adafruit Circuit Playground Express Bluefruit too.   
I'm using a local Mosquitto server running on Raspbian and RPI3.   
Did also some performance tests about Mosquitto and RPI Zero. Had lame results.  

When I restart the Mosquitto server, the client loses the connection both while publishing and subscribing.  
The most effective solution I found is a board reset after a few failed attempts of publishing on an hearth beat Mqtt channel.    