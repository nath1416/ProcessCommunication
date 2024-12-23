#include <stdio.h>
#include <iostream>
#include <sys/types.h>
#include <unistd.h>
#include <iostream>
#include <thread>
#include <chrono>

int main() {
    // std::cout << "UP"<< std::endl;

    std::string line;
    while (std::getline(std::cin, line)) {
        std::cout << "Slave " << line << std::endl;
        // std::cout << "READYd"  << std::endl;
        // break;
        // std::this_thread::sleep_for(std::chrono::seconds(1));
    }

    for(int i = 0; i < 10; i++){
        // std::cout << "PID " << getpid() << std::endl;
    }
    return 0;
}