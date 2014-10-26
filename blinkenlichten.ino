#include <Adafruit_NeoPixel.h>

#define PIN 3
#define NUM_PIXELS 40

// Parameter 1 = number of pixels in strip
// Parameter 2 = Arduino pin number (most are valid)
// Parameter 3 = pixel type flags, add together as needed:
//   NEO_KHZ800  800 KHz bitstream (most NeoPixel products w/WS2812 LEDs)
//   NEO_KHZ400  400 KHz (classic 'v1' (not v2) FLORA pixels, WS2811 drivers)
//   NEO_GRB     Pixels are wired for GRB bitstream (most NeoPixel products)
//   NEO_RGB     Pixels are wired for RGB bitstream (v1 FLORA pixels, not v2)
Adafruit_NeoPixel strip = Adafruit_NeoPixel(NUM_PIXELS, PIN, NEO_GRB + NEO_KHZ800);

// IMPORTANT: To reduce NeoPixel burnout risk, add 1000 uF capacitor across
// pixel power leads, add 300 - 500 Ohm resistor on first pixel's data input
// and minimize distance between Arduino and first pixel.  Avoid connecting
// on a live circuit...if you must, connect GND first.

void setup() {
  Serial.begin(9600);
  Serial.println("ready.");
  strip.begin();
  strip.show(); // Initialize all pixels to 'off'
}

void handleSet() {
  long led;
  long r, g, b;
  led = Serial.parseInt();
  r = Serial.parseInt();
  g = Serial.parseInt();
  b = Serial.parseInt();

  strip.setPixelColor((uint16_t)led, strip.Color(r,g,b));
  strip.show();
}
void handleGet() {
  Serial.println("queries aren't implemented yet. :-(");
}
void handleCount() {
  Serial.println(strip.numPixels());
}

void handleSerial() {
  char next = 0;
  if (Serial.available() > 0) {
    next = Serial.read();
  }

  if (next == 's') {
    handleSet();
  } else if (next == 'g') {
    handleGet();
  } else if (next == 'c') {
    handleCount();
  }
}

void loop() {
  if (Serial.available() > 0) {
    handleSerial();
  }
}

