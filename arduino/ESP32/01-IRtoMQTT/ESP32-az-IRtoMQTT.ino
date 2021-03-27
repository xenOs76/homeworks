/*
 
 Requirements:
 https://github.com/crankyoldgit/IRremoteESP8266
 https://github.com/arduino-libraries/ArduinoMqttClient

*/


#include <Arduino.h>
#include <IRremoteESP8266.h>
#include <IRrecv.h>
#include <IRutils.h>

#include <WiFi.h>
#include <PubSubClient.h>


// Update these with values suitable for your network.
const char* ssid = "CHANGE-ME";
const char* password = "CHANGE-ME";
const char* mqtt_server = "CHANGE-ME";
#define mqtt_port 1883
#define MQTT_USER "CHANGE-ME"
#define MQTT_PASSWORD "CHANGE-ME"
#define MQTT_SERIAL_PUBLISH_CH "/ESP32/ir"
#define MQTT_SERIAL_RECEIVER_CH "/ESP32/ir/rx"


// An IR detector/demodulator is connected to GPIO pin 14(D5 on a NodeMCU
// board).
// Note: GPIO 16 won't work on the ESP8266 as it does not have interrupts.
const uint16_t kRecvPin = 14;
IRrecv irrecv(kRecvPin);
decode_results results;

WiFiClient wifiClient;

PubSubClient client(wifiClient);

void setup_wifi() {
    delay(10);
    // We start by connecting to a WiFi network
    Serial.println();
    Serial.print("Connecting to ");
    Serial.println(ssid);
    WiFi.begin(ssid, password);
    while (WiFi.status() != WL_CONNECTED) {
      delay(500);
      Serial.print(".");
    }
    randomSeed(micros());
    Serial.println("");
    Serial.println("WiFi connected");
    Serial.println("IP address: ");
    Serial.println(WiFi.localIP());
}

void reconnect() {
  // Loop until we're reconnected
  while (!client.connected()) {
    Serial.print("Attempting MQTT connection...");
    // Create a random client ID
    String clientId = "ESP32Client-";
    clientId += String(random(0xffff), HEX);
    // Attempt to connect
    if (client.connect(clientId.c_str(),MQTT_USER,MQTT_PASSWORD)) {
      Serial.println("connected");
      //Once connected, publish an announcement...
      client.publish("/ESP32/status", "online");
    } else {
      Serial.print("failed, rc=");
      Serial.print(client.state());
      Serial.println(" try again in 5 seconds");
      // Wait 5 seconds before retrying
      delay(5000);
    }
  }
}

void callback(char* topic, byte *payload, unsigned int length) {
    Serial.println("-------new message from broker-----");
    Serial.print("channel:");
    Serial.println(topic);
    Serial.print("data:");  
    Serial.write(payload, length);
    Serial.println();
}

void setup() {
  Serial.begin(115200);
  Serial.setTimeout(500);// Set time out for 
  
  irrecv.enableIRIn();  // Start the receiver
  Serial.println();
  Serial.print("IRrecvDemo is now running and waiting for IR message on Pin ");
  Serial.println(kRecvPin);
  
  setup_wifi();
  client.setServer(mqtt_server, mqtt_port);
  client.setCallback(callback);
  reconnect();
}

void publishSerialData(char *serialData){
  if (!client.connected()) {
    reconnect();
  }
  client.publish(MQTT_SERIAL_PUBLISH_CH, serialData);
}
void loop() {
   if (irrecv.decode(&results)) {
    // print() & println() can't handle printing long longs. (uint64_t)
    //publishSerialData(int64String(results.value));
    
    serialPrintUint64(results.value, HEX);
    Serial.println("");

  switch (results.value) {
    case 0xFF6897:
      publishSerialData("1");
      Serial.println("remote: 1");
      break;
    case 0xFF9867:
      publishSerialData("2");
      Serial.println("remote: 2");
      break;
    case 0xFFB04F:
      publishSerialData("3");
      Serial.println("remote: 3");
      break;
    case 0xFF30CF:
      publishSerialData("4");
      Serial.println("remote: 4");
      break;      
    case 0xFF18E7:
      publishSerialData("5");
      Serial.println("remote: 5");
      break; 
    case 0xFF7A85:
      publishSerialData("6");
      Serial.println("remote: 6");
      break;
    case 0xFF10EF:
      publishSerialData("7");
      Serial.println("remote: 7");
      break;
    case 0xFF38C7:
      publishSerialData("8");
      Serial.println("remote: 8");
      break;       
    case 0xFF5AA5:
      publishSerialData("9");
      Serial.println("remote: 9");
      break; 
  }
    
    Serial.println("");
    irrecv.resume();  // Receive the next value
  }
  delay(100);
}
 
