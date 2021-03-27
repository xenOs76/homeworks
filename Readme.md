# Homeworks

>*– Che cosa sia questa musica?*  
>*Peccato che io debba andare a scuola, se no...*  
>*E rimase lì perplesso. A ogni modo, bisognava prendere una risoluzione: o a scuola, o a sentire i pifferi.*  
>*– Oggi anderò a sentire i pifferi, e domani a scuola: per andare a scuola c’è sempre tempo, – disse finalmente quel monello facendo una spallucciata.*  
>
>[Carlo Collodi, Le avventure di Pinocchio](http://www.letteraturaitaliana.net/pdf/Volume_9/t217.pdf)

On my way to becoming what I was meant to be, I got diverted by little things so many times that it became an habit.  
On my way to mastering Python, I got enchanted by [micro](https://micropython.org/) and [circuit](https://circuitpython.org/), the two buddies from the hot board.  
On my way to Go, even then, I stumbled upon some [tiny](https://tinygo.org/) distractions.  
But in any case, I must tell, I had way more fun than expected.  

I'm sharing these **homeworks** because there's always time to get back to school once the music is over.  
Unless, by the way, You're looking for a brand new school book...  

## Changelog

### 20210327 - import Arduino IR receiver and Deep Sleep examples

> *Raccolse poi tutta la paglia che rimaneva all’intorno, e se l’accomodò addosso, facendosene, alla meglio, una specie di coperta, per temperare il freddo, che anche là dentro si faceva sentir molto bene; e vi si rannicchiò sotto, con l’intenzione di dormire un bel sonno, parendogli d’averlo comprato anche più caro del dovere.*
> *Ma appena ebbe chiusi gli occhi, cominciò nella sua memoria o nella sua fantasia (il luogo preciso non ve lo saprei dire), cominciò, dico, un andare e venire di gente, così affollato, così incessante, che addio sonno.*
>
> [Alessandro Manzoni, I Promessi sposi](http://www.letteraturaitaliana.net/pdf/Volume_8/t337.pdf)

### 20210116 - import some Arduino examples

>*– è dolce o amara?*  
>*– è amara, ma ti farà bene.*  
>*– Se è amara, non la voglio.*  
>*– Dà retta a me: bevila.*  
>*– A me l’amaro non mi piace.*  
>*– Bevila: e quando l’avrai bevuta, ti darò una pallina di zucchero, per rifarti la bocca.*  
>*– Dov’è la pallina di zucchero?*  
>*– Eccola qui, – disse la Fata, tirandola fuori da una zuccheriera d’oro.*  
>*– Prima voglio la pallina di zucchero, e poi beverò quell’acquaccia amara...*  

I'm stuck with the [RainbowOnce](tinygo/arduino-nano33/08-mqttSub_NeopixelStrip_RainbowOnce/) example when it comes to receiving commands from Mqtt in a coroutine: *sure, bro, You can spin a rainbow. But You cannot let it shine as long as you want*.  
I feel I'm missing something. I'm pointing in the wrong direction, perhaps.  
Let's get back to [Arduino](https://www.arduino.cc/) IDE and examples while waiting to see my *Nano* standing on the shoulder of giants.  
So I lined up on the breadboard an *[Arduino MKR Wifi 1010](/arduino/ArduinoMKRWifi1010/)* and uploaded some [FastLED](http://fastled.io/) [examples](/arduino/ArduinoMKRWifi1010/00-fastled/), some Mqtt examples and an old classic: the candy from a stranger.  
I uploaded an adapted version of the Arduino sketch from *[Multi-tasking the Arduino](https://learn.adafruit.com/multi-tasking-the-arduino-part-3/overview)*, an [Adafruit](https://www.adafruit.com/) tutorial explaining how to have two or three loops animating LEDs indipendently. Something that should be even easier to do with Tinygo, right?  

### 20210115 - Initial import: examples of Tinygo running on an Arduino Nano33 Iot

I'm looking for a way to control an *[Arduino Nano33 Iot](/tinygo/arduino-nano33/)* over Mqtt.  
I'd like to use the board to get some data from a sensor and drive a Neopixel strip.  
I did translate the Neopixel Rainbow animation from a Python example into Go/Tinygo in the past.  
It used to run on an *Adafruit Circuit Playground Express Bluefruit* too.  
I'm using a local Mosquitto server running on Raspbian and RPI3.  
Did also some performance tests about Mosquitto and RPI Zero. Had lame results.  

When I restart the Mosquitto server, the client loses the connection both while publishing and subscribing.  
The most effective solution I found is a board reset after a few failed attempts of publishing on an hearth beat Mqtt channel.  
