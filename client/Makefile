target/miko.jar:
  mvn package

all: target/miko.jar

start: target/miko.jar
	java -jar target/miko.jar

clean:
	mvn clean

.PHONY: all start clean
