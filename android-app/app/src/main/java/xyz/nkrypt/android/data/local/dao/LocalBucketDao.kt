package xyz.nkrypt.android.data.local.dao

import androidx.room.Dao
import androidx.room.Insert
import androidx.room.OnConflictStrategy
import androidx.room.Query
import kotlinx.coroutines.flow.Flow
import xyz.nkrypt.android.data.local.entity.LocalBucketEntity

@Dao
interface LocalBucketDao {

    @Query("SELECT * FROM local_buckets ORDER BY createdAt DESC")
    fun getAll(): Flow<List<LocalBucketEntity>>

    @Query("SELECT * FROM local_buckets WHERE id = :id")
    suspend fun getById(id: String): LocalBucketEntity?

    @Query("SELECT * FROM local_buckets WHERE name = :name LIMIT 1")
    suspend fun getByName(name: String): LocalBucketEntity?

    @Insert(onConflict = OnConflictStrategy.REPLACE)
    suspend fun insert(entity: LocalBucketEntity)

    @Query("DELETE FROM local_buckets WHERE id = :id")
    suspend fun deleteById(id: String)
}
