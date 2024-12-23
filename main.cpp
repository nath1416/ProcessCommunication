#include <iostream>
#include <thread>
#include <string>
#include <queue>
#include <mutex>
#include <condition_variable>

std::mutex incomingMtx;
std::mutex outcomingMtx;

std::queue<std::string> incomingMessages;
std::queue<std::string> outcomingMessages;

std::condition_variable outcommingCv;
std::condition_variable incommingCv;


template <typename T> 
void addMessageToQueue(std::mutex *mtx, std::condition_variable *cv, std::queue<T> *queue, T *value){
    std::unique_lock<std::mutex> lock(*mtx);
    queue->push(*value);
    lock.unlock();
    cv->notify_one();
}
template <typename T> 
void readMessageFromQueue(std::mutex *mtx, std::condition_variable *cv, std::queue<T> *queue, T *value){
    std::unique_lock<std::mutex> lock(*mtx);
    cv->wait(lock, [queue](){ return !queue->empty(); });
    *value = queue->front();
    queue->pop();
    lock.unlock();
    cv->notify_one();
}

void readFromStdin() {
    std::string line;
    while (true) {
        if (std::getline(std::cin, line)) {
            addMessageToQueue(&incomingMtx, &incommingCv, &incomingMessages, &line);
            addMessageToQueue(&outcomingMtx, &outcommingCv, &outcomingMessages, &line);
        } else {
            break; 
        }
    }
}

void writeToStdout() {
    int count = 0;
    while (true) {
        std::string line;
        readMessageFromQueue(&outcomingMtx, &outcommingCv, &outcomingMessages, &line);
        std::cout << "SLAVE Write to stdout: " << line << std::endl;
    }
}

void changeState() {
    while (true) {
        std::unique_lock<std::mutex> lock(incomingMtx);
        incommingCv.wait(lock, []{return !incomingMessages.empty();});
        std::string message = incomingMessages.front();
        incomingMessages.pop();
        lock.unlock();

    }
}


int main() {
    std::thread reader(readFromStdin);
    std::thread writer(writeToStdout);
    std::thread writers(changeState);

    for (int i = 0; i < 5; ++i) {
        std::this_thread::sleep_for(std::chrono::seconds(1));
    }

    reader.join();
    writer.join();

    return 0;
}