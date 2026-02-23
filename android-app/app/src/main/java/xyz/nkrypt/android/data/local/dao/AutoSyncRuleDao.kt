package xyz.nkrypt.android.data.local.dao

import androidx.room.Dao
import androidx.room.Insert
import androidx.room.OnConflictStrategy
import androidx.room.Query
import kotlinx.coroutines.flow.Flow
import xyz.nkrypt.android.data.local.entity.AutoSyncRuleEntity

@Dao
interface AutoSyncRuleDao {

    @Query("SELECT * FROM auto_sync_rules ORDER BY createdAt DESC")
    fun getAll(): Flow<List<AutoSyncRuleEntity>>

    @Query("SELECT * FROM auto_sync_rules WHERE id = :id")
    suspend fun getById(id: String): AutoSyncRuleEntity?

    @Insert(onConflict = OnConflictStrategy.REPLACE)
    suspend fun insert(entity: AutoSyncRuleEntity)

    @Query("DELETE FROM auto_sync_rules WHERE id = :id")
    suspend fun deleteById(id: String)
}
