# **Lab 6: Producer-Consumer**

This is my solution for the Producer-Consumer lab. It implements a Finite Buffer solution using `sync.Mutex` for exclusive access and two Semaphores (`Spaces` and `Items`) to handle full and empty buffer conditions, preventing race conditions and busy waiting.

## **Setup & Running**

This solution uses the external semaphore package.

1. **Initialize and get dependencies:** Open a terminal in this directory and run:

| go mod init producerconsumer go get golang.org/x/sync/semaphore |
| :---- |

2. **Run the code:**

| go run ProducerConsumer.go |
| :---- |

## **GitHub Repository**

Here is the link to my git repository containing all the labs: [https://github.com/AnJig00/Concurrent-Development-Labs](https://github.com/AnJig00/Concurrent-Development-Labs)

## **License**

This project is licensed under the **GPL v3.0**.

