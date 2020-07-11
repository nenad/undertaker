FROM php:7.2-fpm

RUN apt-get update && apt-get install -y git

RUN git clone https://github.com/krakjoe/tombs.git /tmp/tombs && \
    cd /tmp/tombs && \
    phpize && ./configure && make && make install

RUN php -r "copy('https://getcomposer.org/installer', 'composer-setup.php');" && \
    php composer-setup.php --install-dir=/usr/local/bin --filename=composer && \
    php -r "unlink('composer-setup.php');"

WORKDIR /project

ADD project /project

RUN composer install

ADD tombs.ini /usr/local/etc/php/conf.d/
