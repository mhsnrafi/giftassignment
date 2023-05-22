# Gift Assignment API Readme
Your task is to build an API in Go that can be used to assign gifts to employees. The API should be able to process requests from employees to receive a gift and then assign an available gift based on the employee's interests.

## Requirements
- Go 1.16+
- Gorilla web toolkit v1.8.0

## Installation
1. Clone the repository `git clone https://github.com/mhsnrafi/giftassignment.git`
2. Install the required packages using `go mod download`

## Usage
1. To start the server, run `docker-compose up -d`
2. The server runs on `http://localhost:8080` by default.
3. Available endpoints:
   - `POST /assign-gift/{id}`: Assign a gift to employee to specified by employee id `id`
    
    ```curl --location --request POST 'http://localhost:8080/assign-gift/8' --data-raw```

   - `POST /employees/{employeeId}/return-gift`: Facilitate a gift return for an employee specified by employee id `id`
   ```curl --location --request POST 'http://localhost:8080//employees/1/return-gift' --data-raw```
4. To run the tests, run `go test ./... -v`
5. To stop the server, run `docker-compose down`


## API Documentation

### GiftAssignmentManager
- `AssignGifts(ctx context.Context, employeeID string) (string, error)`: Assigns a gift to the specified employee.
- `ReturnGift(ctx context.Context, employeeID string) error`: Handles the return of a gift from the specified employee.

### Controller
- `AssignGifts(w http.ResponseWriter, r *http.Request)`: Assigns a gift to a user.
- `ReturnGift(w http.ResponseWriter, r *http.Request)`: Handles the return of a gift from a user


### EmployeeGiftManagementClient.AssignGifts and ReturnGift Function Explanation
The EmployeeGiftManagementClient struct contains methods AssignGifts and ReturnGift to manage gifts for employees, ensuring the entire process is carried out smoothly and concurrently. It interacts with the underlying storage client to access and manipulate data. Here's a detailed explanation of how these methods work:

NewEmployeeGiftManagementClient
This function acts as a constructor, returning a new instance of EmployeeGiftManagementClient with a provided storage client and locker for distributed lock management. This setup allows the client to execute operations on the underlying storage system and manage concurrency with the locker.

### AssignGifts
The AssignGifts method is responsible for assigning a suitable gift to a specific employee. It employs the concept of 'operation', a function passed to withLockAndTransaction that performs the assignment under a lock and transaction, ensuring data consistency and atomicity.

Here's how AssignGifts works:

1.` Validate employee ID`: This checks if the given employee ID is valid. If not, it logs an error and returns a custom error ErrInvalidEmployeeID.

2. `Get employee's details and verify if a gift is already assigned`: If the employee already has a gift that hasn't been returned yet, an error ErrGiftAlreadyAssigned is returned.

3. `Check available gifts`: All available gifts are fetched from the repository. If this operation fails, an error is returned.

4. `Assign a suitable gift`: The method iterates over each available gift and the employee's interests. If it finds a matching category, the gift is assigned to the employee, and both gift and employee records are updated. If all operations are successful, it returns the name of the assigned gift.

5. `Assign a random gift if no suitable gift is found`: If no matching gift is found, a random gift from the available ones is assigned to the employee, and the respective records are updated.


### ReturnGift
The ReturnGift method allows an employee to return a previously assigned gift. It follows the same 'operation' logic, executing the return process under a lock and transaction.

Here's how ReturnGift works:

1. `Validate employee ID`: This checks if the given employee ID is valid. If not, it logs an error and returns a custom error ErrInvalidEmployeeID.

2. `Get the employee's gift details`: The details of the employee's assigned gift are fetched.

3. `Verify if the gift can be returned`: If the employee doesn't have a gift or has already returned a gift, a custom error ErrGiftNotReturnable is returned.

4. `Mark the gift as returned`: If the gift can be returned, both the gift and employee records are updated to mark the gift as returned.

The entire return process is performed under a lock, provided by withLockAndTransaction, ensuring safe concurrent operations.

### withLockAndTransaction
This method is a helper function that enables executing the provided operation under a distributed lock and database transaction. It performs the following steps:

1. `Acquire the distributed lock`: It attempts to acquire the distributed lock for a given employeeID. If it fails to do so, it logs an error and returns it.

2. `Start a new transaction`: It starts a new transaction using EmployeeRepo.BeginTransaction. If it fails, it logs an error and returns it.

3.`Perform the operation:` It executes the operation function, passing the transaction context and transaction object.

4. `Commit the transaction`: If the operation is successful, it commits the transaction, making all changes permanent in the database.

5. `Handle errors and rollbacks`: If any errors occur during the operation or while committing the transaction, it rolls back the changes to maintain database consistency.

6. `Release the lock`: Finally, it releases the distributed lock, allowing other operations to proceed.

The method ensures data consistency and atomicity by using locks and transactions, facilitating safe concurrent operations.

## Answer to questions:

## **A gift can only be given to one employee, and we expect all employees to rush to this website. What are potential issues with this rush?**
Potential issues include server overload, network latency, and race conditions. Server overload and network latency occur when a large number of users access the website at the same time, leading to performance issues and potential crashes. Race conditions can occur when multiple employees try to claim the same gift at the same time. This can be handled using locks to ensure that only one operation happens per employee at a time. The code now includes a distributed locking mechanism which allows it to handle a high number of concurrent requests. By locking on the employee ID, we make sure that each employee can only get a gift once, even if they try to get a gift multiple times in quick succession
# **It’s possible we might not have bought appropriate gifts, how can we try to increase the number of employees who receive gifts?**
To increase the number of employees who receive gifts, we first attempt to match gifts based on employee interests (assignGiftBasedOnInterest). If no suitable gift is found based on interests, we assign a random gift to the employee (assignRandomGift).


# **Once an employee has been assigned a gift that’s it, no changing. How can we enforce this**
The enforcement of no gift changes is handled by checking if the employee already has a gift that hasn't been returned (getValidEmployeeWithoutGift). If the employee already has a gift that hasn't been returned, an error is returned

# **In the future we might want employees to be able to return a gift once, think about (but don’t implement!) how this might work with your solution.**
If gift returns are allowed in the future, we can implement a ReturnGift function similar to the one in the provided code. The function would check if the employee is valid and has a gift to return. If so, it would process the gift return and mark the gift as returned in the gift and employee records

# **How would this work at “Santa scale”? For your takehome we advise you to K.I.S.S., but be prepared to chat about limitations/tradeoffs in your approach**
Handling this gift assignment system at a "Santa scale" implies managing a large number of employees, possibly in the millions or more. At such a scale, several factors must be considered, and there are trade-offs associated with each approach.

Concurrency and Load Balancing: In the current implementation, we are using distributed locks to handle concurrency issues, which could become a bottleneck at large scale due to lock contention. An alternative approach could be optimistic concurrency control where we don't lock at the beginning but rather check during update or commit to see if the record was updated by another transaction. This reduces lock contention but might increase the number of transaction retries.

Database Scaling: The current system likely assumes a single database instance, which could become a bottleneck at high scale. To handle larger loads, we might need to implement database sharding, where data is partitioned across multiple databases. This allows for horizontal scaling and increased capacity, but adds complexity to the database architecture and application logic. Another approach could be using a NoSQL database designed for high-scale read and write operations.

Caching: Implementing caching could be another way to improve performance. We could cache the available gifts in memory to speed up gift assignments. However, this approach requires a good cache invalidation strategy and can introduce consistency issues.



