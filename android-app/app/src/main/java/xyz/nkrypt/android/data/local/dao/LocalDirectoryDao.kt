package xyz.nkrypt.android.data.local.dao

import androidx.room.Dao
import androidx.room.Insert
import androidx.room.OnConflictStrategy
import androidx.room.Query
import xyz.nkrypt.android.data.local.entity.LocalDirectoryEntity

@Dao
interface LocalDirectoryDao {

    @Query("SELECT * FROM local_directories WHERE bucketId = :bucketId AND (parentDirectoryId = :parentId OR (:parentId IS NULL AND parentDirectoryId IS NULL)) ORDER BY name")
    suspend fun getByBucketAndParent(bucketId: String, parentId: String?): List<LocalDirectoryEntity>

    @Query("SELECT * FROM local_directories WHERE id = :id")
    suspend fun getById(id: String): LocalDirectoryEntity?

    @Insert(onConflict = OnConflictStrategy.REPLACE)
    suspend fun insert(entity: LocalDirectoryEntity)

    @Query("UPDATE local_directories SET name = :newName WHERE id = :id")
    suspend fun rename(id: String, newName: String)

    @Query("UPDATE local_directories SET parentDirectoryId = :newParentId WHERE id = :id")
    suspend fun move(id: String, newParentId: String?)

    @Query("DELETE FROM local_directories WHERE id = :id")
    suspend fun deleteById(id: String)
}
