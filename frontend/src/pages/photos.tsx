import React, { useState } from 'react';
import Head from 'next/head';
import { Grid, List, Plus, Upload } from 'lucide-react';
import { Navbar } from '@/components/facebook/Navbar';
import { Button } from '@/components/ui/button';
import { Card, CardContent } from '@/components/ui/card';
import { Dialog, DialogContent, DialogHeader, DialogTitle, DialogTrigger } from '@/components/ui/dialog';

const PhotosPage: React.FC = () => {
  const [viewMode, setViewMode] = useState<'grid' | 'list'>('grid');
  const [selectedPhoto, setSelectedPhoto] = useState<string | null>(null);

  const albums = [
    {
      id: '1',
      title: 'Profile Pictures',
      photoCount: 24,
      coverPhoto: 'https://picsum.photos/seed/album1/600/400',
      lastUpdated: '2 days ago',
    },
    {
      id: '2',
      title: 'Cover Photos',
      photoCount: 12,
      coverPhoto: 'https://picsum.photos/seed/album2/600/400',
      lastUpdated: '1 week ago',
    },
    {
      id: '3',
      title: 'Vacation 2024',
      photoCount: 156,
      coverPhoto: 'https://picsum.photos/seed/album3/600/400',
      lastUpdated: '3 weeks ago',
    },
    {
      id: '4',
      title: 'Friends & Family',
      photoCount: 89,
      coverPhoto: 'https://picsum.photos/seed/album4/600/400',
      lastUpdated: '1 month ago',
    },
  ];

  const recentPhotos = Array.from({ length: 24 }, (_, i) => ({
    id: i + 1,
    url: `https://picsum.photos/seed/photo${i}/600/400`,
    caption: `Photo ${i + 1}`,
    likes: Math.floor(Math.random() * 100),
    comments: Math.floor(Math.random() * 20),
  }));

  return (
    <>
      <Head>
        <title>Photos - Facebook Clone</title>
        <meta name="description" content="View and manage your photos" />
      </Head>

      <div className="min-h-screen bg-gray-100 dark:bg-gray-950">
        <Navbar />

        <div className="max-w-7xl mx-auto p-4">
          {/* Header */}
          <div className="flex items-center justify-between mb-6">
            <h1 className="text-3xl font-bold">Photos</h1>
            <div className="flex space-x-2">
              <Dialog>
                <DialogTrigger asChild>
                  <Button variant="default">
                    <Upload className="w-4 h-4 mr-2" />
                    Upload Photos
                  </Button>
                </DialogTrigger>
                <DialogContent>
                  <DialogHeader>
                    <DialogTitle>Upload Photos</DialogTitle>
                  </DialogHeader>
                  <div className="border-2 border-dashed border-gray-300 dark:border-gray-700 rounded-lg p-12 text-center">
                    <Upload className="w-12 h-12 mx-auto mb-4 text-gray-400" />
                    <p className="text-gray-600 dark:text-gray-400">
                      Drag and drop photos here or click to browse
                    </p>
                  </div>
                </DialogContent>
              </Dialog>
              <Button
                variant="outline"
                size="icon"
                onClick={() => setViewMode('grid')}
                className={viewMode === 'grid' ? 'bg-blue-50 dark:bg-blue-950' : ''}
              >
                <Grid className="w-4 h-4" />
              </Button>
              <Button
                variant="outline"
                size="icon"
                onClick={() => setViewMode('list')}
                className={viewMode === 'list' ? 'bg-blue-50 dark:bg-blue-950' : ''}
              >
                <List className="w-4 h-4" />
              </Button>
            </div>
          </div>

          {/* Albums Section */}
          <div className="mb-8">
            <h2 className="text-2xl font-bold mb-4">Albums</h2>
            <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-4">
              {/* Create Album Card */}
              <Card className="cursor-pointer hover:shadow-lg transition-shadow">
                <CardContent className="p-0">
                  <div className="aspect-[4/3] bg-gray-200 dark:bg-gray-800 flex flex-col items-center justify-center">
                    <div className="w-12 h-12 bg-gray-300 dark:bg-gray-700 rounded-full flex items-center justify-center mb-2">
                      <Plus className="w-6 h-6" />
                    </div>
                    <p className="font-semibold">Create Album</p>
                  </div>
                </CardContent>
              </Card>

              {/* Album Cards */}
              {albums.map((album) => (
                <Card key={album.id} className="cursor-pointer hover:shadow-lg transition-shadow">
                  <CardContent className="p-0">
                    <div className="aspect-[4/3] overflow-hidden">
                      <img
                        src={album.coverPhoto}
                        alt={album.title}
                        className="w-full h-full object-cover hover:scale-105 transition-transform duration-200"
                      />
                    </div>
                    <div className="p-4">
                      <h3 className="font-semibold text-lg mb-1">{album.title}</h3>
                      <p className="text-sm text-gray-600 dark:text-gray-400">
                        {album.photoCount} photos · {album.lastUpdated}
                      </p>
                    </div>
                  </CardContent>
                </Card>
              ))}
            </div>
          </div>

          {/* Recent Photos */}
          <div>
            <h2 className="text-2xl font-bold mb-4">Recent Photos</h2>
            {viewMode === 'grid' ? (
              <div className="grid grid-cols-2 md:grid-cols-3 lg:grid-cols-4 xl:grid-cols-6 gap-2">
                {recentPhotos.map((photo) => (
                  <Dialog key={photo.id}>
                    <DialogTrigger asChild>
                      <div className="aspect-square overflow-hidden rounded-lg cursor-pointer group">
                        <img
                          src={photo.url}
                          alt={photo.caption}
                          className="w-full h-full object-cover hover:scale-110 transition-transform duration-200"
                        />
                      </div>
                    </DialogTrigger>
                    <DialogContent className="max-w-4xl">
                      <div className="flex items-center justify-center">
                        <img
                          src={photo.url}
                          alt={photo.caption}
                          className="max-h-[80vh] object-contain"
                        />
                      </div>
                      <div className="mt-4">
                        <p className="font-semibold">{photo.caption}</p>
                        <p className="text-sm text-gray-600 dark:text-gray-400">
                          {photo.likes} likes · {photo.comments} comments
                        </p>
                      </div>
                    </DialogContent>
                  </Dialog>
                ))}
              </div>
            ) : (
              <div className="space-y-4">
                {recentPhotos.map((photo) => (
                  <Card key={photo.id}>
                    <CardContent className="p-4 flex items-center space-x-4">
                      <div className="w-32 h-32 overflow-hidden rounded-lg flex-shrink-0">
                        <img
                          src={photo.url}
                          alt={photo.caption}
                          className="w-full h-full object-cover"
                        />
                      </div>
                      <div className="flex-1">
                        <h3 className="font-semibold text-lg mb-2">{photo.caption}</h3>
                        <p className="text-sm text-gray-600 dark:text-gray-400">
                          {photo.likes} likes · {photo.comments} comments
                        </p>
                      </div>
                      <Button variant="outline">View</Button>
                    </CardContent>
                  </Card>
                ))}
              </div>
            )}
          </div>
        </div>
      </div>
    </>
  );
};

export default PhotosPage;
