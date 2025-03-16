package taskService

import (
     "encoding/json"
     "fmt"
     "time"
    )

    type Tasks struct {
     ID          uint      `json:"id"`
     Title       string    `json:"title"`
     Description string    `json:"description"`
     CompletedAt time.Time `json:"completed_at"`
     IsDone      bool `json:"is_done"`  
    }

    func (m *Tasks) UnmarshalJSON(data []byte) error {
     type Alias Tasks

     aux := &struct {
      CompletedAt string `json:"completed_at"`
      *Alias
     }{
      Alias: (*Alias)(m),
     }

     if err := json.Unmarshal(data, &aux); err != nil {
      return err
     }

     if aux.CompletedAt == "completed" {
      m.CompletedAt = time.Time{}
      return nil
     }

     if aux.CompletedAt == "" {
      m.CompletedAt = time.Time{}
      return nil
     }

     layout := "2006-01-02 15:04:05 +0000 UTC"
     fmt.Printf("Unmarshaling completed_at: %s\n", aux.CompletedAt)
     t, err := time.Parse(layout, aux.CompletedAt)
     if err != nil {
      return fmt.Errorf("invalid date format: %w", err)
     }
     m.CompletedAt = t
     return nil
    }