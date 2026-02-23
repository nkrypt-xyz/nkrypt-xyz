package xyz.nkrypt.android.di

import android.content.Context
import androidx.room.Room
import dagger.Module
import dagger.Provides
import dagger.hilt.InstallIn
import dagger.hilt.android.qualifiers.ApplicationContext
import dagger.hilt.components.SingletonComponent
import xyz.nkrypt.android.data.local.AppDatabase
import xyz.nkrypt.android.data.local.dao.AutoImportRuleDao
import xyz.nkrypt.android.data.local.dao.AutoSyncRuleDao
import xyz.nkrypt.android.data.local.dao.LocalBlobDao
import xyz.nkrypt.android.data.local.dao.LocalBucketDao
import xyz.nkrypt.android.data.local.dao.LocalDirectoryDao
import xyz.nkrypt.android.data.local.dao.LocalFileDao
import xyz.nkrypt.android.data.local.dao.RemoteBucketDao
import javax.inject.Singleton

@Module
@InstallIn(SingletonComponent::class)
object AppModule {

    @Provides
    @Singleton
    fun provideContext(@ApplicationContext context: Context): Context = context

    @Provides
    @Singleton
    fun provideAppDatabase(@ApplicationContext context: Context): AppDatabase =
        Room.databaseBuilder(context, AppDatabase::class.java, "nkrypt.db")
            .fallbackToDestructiveMigration()
            .build()

    @Provides
    @Singleton
    fun provideLocalBucketDao(db: AppDatabase): LocalBucketDao = db.localBucketDao()

    @Provides
    @Singleton
    fun provideLocalDirectoryDao(db: AppDatabase): LocalDirectoryDao = db.localDirectoryDao()

    @Provides
    @Singleton
    fun provideLocalFileDao(db: AppDatabase): LocalFileDao = db.localFileDao()

    @Provides
    @Singleton
    fun provideLocalBlobDao(db: AppDatabase): LocalBlobDao = db.localBlobDao()

    @Provides
    @Singleton
    fun provideRemoteBucketDao(db: AppDatabase): RemoteBucketDao = db.remoteBucketDao()

    @Provides
    @Singleton
    fun provideAutoImportRuleDao(db: AppDatabase): AutoImportRuleDao = db.autoImportRuleDao()

    @Provides
    @Singleton
    fun provideAutoSyncRuleDao(db: AppDatabase): AutoSyncRuleDao = db.autoSyncRuleDao()
}
