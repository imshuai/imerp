package main

import (
	"fmt"
	"erp/models"
)

func main() {
	fmt.Println("TaskStatusPending:", string(models.TaskStatusPending))
	fmt.Println("TaskStatusInProgress:", string(models.TaskStatusInProgress))
	fmt.Println("TaskStatusCompleted:", string(models.TaskStatusCompleted))
}
