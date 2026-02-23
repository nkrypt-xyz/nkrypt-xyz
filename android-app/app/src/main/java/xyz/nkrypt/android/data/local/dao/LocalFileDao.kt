package xyz.nkrypt.android.data.local.dao

import androidx.room.Dao
import androidx.room.Insert
import androidx.room.OnConflictStrategy
import androidx.room.Query
import xyz.nkrypt.android.data.local.entity.LocalFileEntity

@Dao
interface LocalFileDao {

    @Query("SELECT * FROM local_files WHERE bucketId = :bucketId AND (directoryId = :dirId OR (:dirId IS NULL AND directoryId IS NULL)) ORDER BY name")
    suspend fun getByBucketAndDirectory(bucketId: String, dirId: String?): List<LocalFileEntity>

    @Query("SELECT * FROM local_files WHERE bucketId = :bucketId AND name = :name AND (directoryId = :dirId OR (:dirId IS NULL AND directoryId IS NULL)) LIMIT 1")
    suspend fun getByNameAndDirectory(bucketId: String, dirId: String?, name: String): LocalFileEntity?

    @Query("SELECT * FROM local_files WHERE id = :id")
    suspend fun getById(id: String): LocalFileEntity?

    @Insert(onConflict = OnConflictStrategy.REPLACE)
    suspend fun insert(entity: LocalFileEntity)

    @Query("UPDATE local_files SET name = :newName WHERE id = :id")
    suspend fun rename(id: String, newName: String)

    @Query("UPDATE local_files SET directoryId = :newDirId WHERE id = :id")
    suspend fun move(id: String, newDirId: String?)

    @Query("DELETE FROM local_files WHERE id = :id")
    suspend fun deleteById(id: String)
}
