<?php

namespace App\Controller;

use Doctrine\DBAL\Connection;
use Symfony\Component\HttpFoundation\JsonResponse;
use Symfony\Component\Routing\Annotation\Route;

class HealthController
{
    private Connection $connection;

    // Внедрение зависимости через конструктор
    public function __construct(Connection $connection)
    {
        $this->connection = $connection;
    }

    #[Route('/api/health', name: 'health_check')]
    public function check(): JsonResponse
    {
        try {
            $this->connection->connect();
            return new JsonResponse(['status' => 'Database is accessible']);
        } catch (\Exception $e) {
            return new JsonResponse(['status' => 'Database is not accessible', 'error' => $e->getMessage()], 500);
        }
    }
}
