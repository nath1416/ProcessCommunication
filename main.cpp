#include <stdio.h>
#include <iostream>
#include <sys/types.h>
#include <unistd.h>

int main() {
    // std::cout << "UP"<< std::endl;

    std::string line;
    while (std::getline(std::cin, line)) {
        std::cout << "PID " << line << std::endl;
        // std::cout << "READYd"  << std::endl;
        // break;
    }

    for(int i = 0; i < 10; i++){
        // std::cout << "PID " << getpid() << std::endl;
    }
    return 0;
}