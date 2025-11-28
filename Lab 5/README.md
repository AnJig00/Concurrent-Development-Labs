# **Lab 5: Dining Philosophers**

This is my solution for the Dining Philosophers lab. The original template forces a deadlock where everyone picks up their left fork simultaneously. I fixed this using the **Asymmetry Solution** (from the slides), where one philosopher picks up their forks in reverse order (Right then Left), making deadlock impossible.

## **Setup & Running**

This solution uses standard Go libraries, so no external installation is required.

1. Run the code(**Open a terminal in this directory and run):

```bash
go run dinPhil.go
```
## **GitHub Repository**

Here is the link to my git repository containing all the labs: [https://github.com/AnJig00/Concurrent-Development-Labs](https://github.com/AnJig00/Concurrent-Development-Labs)

## **License**

This project is licensed under the **GPL v3.0**.

