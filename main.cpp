#include <iostream>
#include <thread>
#include <string>
#include <queue>
#include <mutex>

std::mutex incomingMtx;
std::mutex outcomingMtx;

std::queue<std::string> incomingMessages;
std::queue<std::string> outcomingMessages;


void readFromStdin() {
    std::string line;
    while (true) {
        if (std::getline(std::cin, line)) {
            {
            std::lock_guard<std::mutex> lock(incomingMtx);
            incomingMessages.push(line);
            }
            {
            std::lock_guard<std::mutex> lock(outcomingMtx);
                outcomingMessages.push("received a packet");
            }
            // std::cout << "SLAVE Read from stdin: " << line << std::endl;

        } else {
            break; 
        }
    }
}

void writeToStdout() {
    int count = 0;
    while (true) {
        while(!outcomingMessages.empty()){
        {
            std::lock_guard<std::mutex> lock(outcomingMtx);
            std::string message = outcomingMessages.front();
            outcomingMessages.pop();
            std::cout << "SLAVE Write to stdout: " << message << std::endl;
        }
            // std::this_thread::sleep_for(std::chrono::milliseconds(100)); 
        }
    }
}

int main() {
    std::thread reader(readFromStdin);
    std::thread writer(writeToStdout);

    // std::cout << "Main thread is running..." << std::endl;
    for (int i = 0; i < 5; ++i) {
        std::this_thread::sleep_for(std::chrono::seconds(1));
    }

    reader.join();
    writer.join();

    return 0;
}