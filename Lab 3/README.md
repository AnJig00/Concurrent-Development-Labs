# **Lab 3: Barrier**

This is my solution for the Barrier synchronization lab. It uses a `sync.Mutex` and the weighted semaphore package (`golang.org/x/sync/semaphore`) to ensure all routines finish Part A before any of them proceed to Part B.

## **Setup & Running**

Since this code uses an external semaphore package, you need to download the dependencies first.

1. **Initialize and get dependencies:** Open a terminal in this directory and run:

```bash
   go mod init barrier
   go get golang.org/x/sync/semaphore
```
2. **Run the code:**

```bash
go run Barrier.go
```
## **GitHub Repository**

Here is the link to my git repository containing all the labs: [https://github.com/AnJig00/Concurrent-Development-Labs](https://github.com/AnJig00/Concurrent-Development-Labs)

## **License**

This project is licensed under the **GPL v3.0**.

