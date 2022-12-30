import { Image } from 'antd';
import type { FC  } from 'react';
import type { ImageProps  } from 'antd';
import { memo } from 'react';
import ImageError from '@/images/image-error.png';

/**
 * encapsulation Image component of the antd
 * 1. set default fallback of Image
 */
const NewImage: FC<ImageProps> = memo((props) => {
  return <Image fallback={ImageError} {...props} />;
});

export default NewImage;
