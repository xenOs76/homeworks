/*
 * https://www.arduino.cc/en/Guide/MKRWiFi1010/powering-with-batteries
 * https://www.arduino.cc/en/Guide/MKRWiFi1010/connecting-sensors
 * https://arduinojson.org/v6/how-to/use-arduinojson-with-arduinomqttclient/
 * https://arduino-tutorials.net/tutorial/capacitive-soil-moisture-sensor-arduino
 * https://arduinojson.org/v6/how-to/use-arduinojson-with-arduinomqttclient/
 */

#include <ArduinoMqttClient.h>
#include <WiFiNINA.h>
#include <ArduinoLowPower.h>
#include <BME280I2C.h>
#include <Wire.h>
#include <ArduinoJson.h>

int resetPin = 3;

// WIFI
#include "arduino_secrets.h"
///////please enter your sensitive data in the Secret tab/arduino_secrets.h
char ssid[] = SECRET_SSID;   // your network SSID (name)
char pass[] = SECRET_PASS;   // your network password (use for WPA, or use as key for WEP)
int keyIndex = 0;            // your network key Index number (needed only for WEP)
int status = WL_IDLE_STATUS;
WiFiClient client;  // Initialize the Wifi client library
WiFiClient wifiClient;

// MQTT
MqttClient mqttClient(wifiClient);
const char broker[] = "192.168.0.24";
int        port     = 1883;
const char topic[]  = "arduino/mkr1010";

#define ms_per_min    60000
#define ms_per_sec    1000
const unsigned long lowPowerSleepInterval = 10 * ms_per_min;

// moisture
const int dry = 912; // value for dry sensor
const int wet = 385; // value for wet sensor

// temp 
BME280I2C bme;    // Default : forced mode, standby time = 1000 ms
                  // Oversampling = pressure ×1, temperature ×1, humidity ×1, filter off,

void setup() {

  digitalWrite(resetPin, HIGH);
  delay(200);
  pinMode(resetPin, OUTPUT);   

  Serial.begin(9600);

  if (WiFi.status() == WL_NO_MODULE) {
    Serial.println("Communication with WiFi module failed!");
    while (true);
  }

  while (status != WL_CONNECTED) {
    Serial.print("Attempting to connect to SSID: ");
    Serial.println(ssid);
    status = WiFi.begin(ssid, pass);
    WiFi.lowPowerMode();
    delay(5000);
  }

  Serial.print("Attempting to connect to the MQTT broker: ");
  Serial.println(broker);

  if (!mqttClient.connect(broker, port)) {
    Serial.print("MQTT connection failed! Error code = ");
    Serial.println(mqttClient.connectError());
    while (1);
  }

  Wire.begin();
  while(!bme.begin())
  {
    Serial.println("Could not find BME280 sensor!");
    delay(1000);
  }

  switch(bme.chipModel())
  {
     case BME280::ChipModel_BME280:
       Serial.println("Found BME280 sensor! Success.");
       break;
     case BME280::ChipModel_BMP280:
       Serial.println("Found BMP280 sensor! No Humidity available.");
       break;
     default:
       Serial.println("Found UNKNOWN sensor! Error!");
  }
}




void loop() {
 
  // moisture
  int moistVal = analogRead(A0);
  int percMoistHum = map(moistVal, wet, dry, 100, 0);

  StaticJsonDocument<200> moist;
  moist["sensor"] = "moistSens2.0";
  moist["location"] = "succulent";
  moist["perc"] = percMoistHum;
  serializeJson(moist, Serial);
  Serial.println("");
 
  mqttClient.beginMessage("moist");
  serializeJson(moist, mqttClient);
  mqttClient.endMessage();
  
  // temp, hum, press
  float temp(NAN), hum(NAN), pres(NAN);
  BME280::TempUnit tempUnit(BME280::TempUnit_Celsius);
  BME280::PresUnit presUnit(BME280::PresUnit_Pa);
  bme.read(pres, temp, hum, tempUnit, presUnit);

  StaticJsonDocument<200> doc;
  doc["sensor"]= "bme280";
  doc["location"]= "indoor";
  doc["temp"] = temp;
  doc["hum"] = hum;
  doc["pres"] = pres;
  serializeJson(doc, Serial);
  Serial.println();
  
  mqttClient.beginMessage("temp");
  serializeJson(doc, mqttClient);
  mqttClient.endMessage();
  delay(300);
  
  WiFi.end();
  Serial.println("going to sleep");
  Serial.end();
  LowPower.deepSleep(lowPowerSleepInterval);
  digitalWrite(resetPin, LOW);

}
