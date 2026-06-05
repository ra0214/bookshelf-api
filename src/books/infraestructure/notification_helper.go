package infraestructure

import "log"

// NotifyTaskCreated envía notificación cuando se crea una nueva tarea
// TODO: Implementar integración con servicio de notificaciones (FCM, WebPush, etc.)
func NotifyTaskCreated(taskID int32, taskTitle string, userID int32) {
	log.Printf("[Notifications] Nueva tarea creada - ID: %d, Título: %s, Usuario: %d", taskID, taskTitle, userID)
}

// NotifyTaskUpdated envía notificación cuando se actualiza una tarea
// TODO: Implementar integración con servicio de notificaciones
func NotifyTaskUpdated(taskID int32, taskTitle string, userID int32) {
	log.Printf("[Notifications] Tarea actualizada - ID: %d, Título: %s, Usuario: %d", taskID, taskTitle, userID)
}

// NotifyTaskDeleted envía notificación cuando se elimina una tarea
// TODO: Implementar integración con servicio de notificaciones
func NotifyTaskDeleted(taskID int32, userID int32) {
	log.Printf("[Notifications] Tarea eliminada - ID: %d, Usuario: %d", taskID, userID)
}
