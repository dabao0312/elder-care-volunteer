package tasks

import (
	"log"
	"time"

	"gorm.io/gorm"

	"elder-care-volunteer/models"
	"elder-care-volunteer/services"
)

func StartNoReplyChecker(db *gorm.DB) {

	go func() {

		for {
			log.Println("[TASK] checking elders no reply...")

			var elders []models.Elder
			db.Find(&elders)

			for _, elder := range elders {

				// ⚠️ 这里先用 CreatedAt 模拟“最后回复时间”
				lastReply := elder.CreatedAt

				if time.Since(lastReply) > 2*time.Hour {

					// log.Printf(
					// 	"[TASK] elder %d exceed 2h, should trigger alert",
					// 	elder.ID,
					// )

					services.AlertNearestVolunteer(db, elder)

					// ⚠️ 这里只是打印
					// 下一步我们会真正调用匹配逻辑
				}
			}

			time.Sleep(1 * time.Minute)
		}
	}()
}
