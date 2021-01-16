#include <FastLED.h>
#define NUM_LEDS 60

CRGBArray<NUM_LEDS> leds;

void setup() { FastLED.addLeds<NEOPIXEL,5>(leds, NUM_LEDS); }

void loop(){ 
  static uint8_t hue;
  for(int i = 0; i < NUM_LEDS/2; i++) {   
    // fade everything out
    leds.fadeToBlackBy(60);

    // let's set an led value
    leds[i] = CHSV(hue++,255,255);

    // now, let's first 20 leds to the top 20 leds, 
    leds(NUM_LEDS/2,NUM_LEDS-1) = leds(NUM_LEDS/2 - 1 ,0);
    FastLED.delay(33);
  }
}
